#!/bin/bash
# This script can be called periodically to sync photographs between
# Phonebook server instances.

#===================================================
#  UPLOAD PORTION
#  If anything has changed in the picture directory
#  then update the repo
#===================================================
WORKDIR="picsync"
STATFILE="picstat"
INSTANCES="pbinst.txt"
UPLOAD=1
MYNAME=$(curl -s http://169.254.169.254/latest/meta-data/public-hostname)
ARTIFPATH="accord/db/picsync"
DEPLOY="/usr/local/accord/bin/deployfile.sh"
GETFILE="/usr/local/accord/bin/getfile.sh"

if [ ! -d "${WORKDIR}" ]; then
	mkdir "${WORKDIR}"
fi

ls -l pictures | cksum > ${WORKDIR}/${STATFILE}new

#===================================
#  From phonebook dir to ${WORKDIR}
#===================================
pushd "${WORKDIR}"
if [ -e "${STATFILE}" ]; then
	$(diff ${STATFILE} ${STATFILE}new > x)
	if [ 0 -eq $(cat x | wc -c) ]; then
		UPLOAD=0
	fi
else
	mv "${STATFILE}new" "${STATFILE}"
fi

if [ ${UPLOAD} -eq 1 ]; then
	cd ..;tar cf "${WORKDIR}"/pictures.tar pictures;cd ${WORKDIR}
	gzip pictures.tar
	${DEPLOY} ${STATFILE} "${ARTIFPATH}/${MYNAME}" >/dev/null
	${DEPLOY} pictures.tar.gz "${ARTIFPATH}/${MYNAME}" >/dev/null
	echo "************************************************"
	echo "***  uploaded ${STATFILE} and pictures.tar.gz" 
	echo "************************************************"
	rm pictures.tar.gz
else
	echo "UPLOAD not needed"
fi

if [ -f "${STATFILE}new" ]; then
	mv ${STATFILE}new ${STATFILE}
fi

rm -f x

#===================================================
#  DOWNLOAD PORTION
#  Get the list of instances, and sync
#===================================================
${GETFILE} ${ARTIFPATH}/${INSTANCES}

cat ${INSTANCES} | while read inst || [[ -n "$inst" ]];
do
	if [ ${inst} != ${MYNAME} ]; then
		echo "Process ${inst}"
		if [ ! -d ${inst} ]; then
			mkdir ${inst}
		fi
		#==============================
		#  From ${WORKDIR} TO ${INST}
		#==============================
		pushd ${inst}
		if [ ! -f ${STATFILE} ]; then
			${GETFILE} "${ARTIFPATH}/${inst}/${STATFILE}"
		fi
		DOWNLOAD=1
		#---------------------------------------------------------------------
		# if there's still no statfile, then there's nothing to download
		# it means that this node has not yet made its first file publish
		#---------------------------------------------------------------------
		if [ ! -f ${STATFILE} ]; then
			DOWNLOAD=0
		elif [ -f ${STATFILE}prior ]; then
			$(diff ${STATFILE} ${STATFILE}prior > x)
			if [ 0 -eq $(cat x | wc -c) ]; then
				DOWNLOAD=0
			fi
		fi
		if [ -f ${STATFILE} ]; then
			mv "${STATFILE}" "${STATFILE}prior" 
		fi
		if [ ${DOWNLOAD} -eq 1 ]; then
			# download this instance's pictures
			${GETFILE} "${ARTIFPATH}/${inst}/pictures.tar.gz"
			gunzip pictures.tar.gz 
			tar xvf pictures.tar
			rsync -auv pictures/* ../../pictures/
			rm -rf pictures.tar pictures
			#============================
			#  From ${INST} TO PICTURES/
			#============================
			#-------------------------------------------------------------------
			# Three step approach to find files of the same name but 
			# different extensions.
			#-------------------------------------------------------------------
			pushd ../../pictures/
			arr=()
			ext=()
			fmatch=()
			ematch1=()
			ematch2=()
			#-------------------------------------------------------------------
			# Step 1: load all files and their extensions into separate arrays
			#-------------------------------------------------------------------
			for f in *; do
		        arr+=(${f%.*})       ## save file name portion
		        ext+=(${f##*.})      ## save extension portion
			done

			len=${#arr[@]}
			# echo "arr has ${len} items"

			#-----------------------------------------------------------
			# Step 2: look for 2 consecutive identical file names
			#-----------------------------------------------------------
			if [ ${len} -gt 1 ]; then
		        tLen=$((${len}-1))
		        echo "loop for i = 0; i < ${tLen}"
		        for (( i=0; i<${tLen}; i++ )); do
	                if [ ${arr[$i]} == ${arr[$((i+1))]} ]; then
                        echo "FOUND at $i, ${arr[$i]}"
                        fmatch+=(${arr[$i]})
                        ematch1+=(${ext[$i]})
                        ematch2+=(${ext[$((i+1))]})
	                fi  
		        done
			fi

			#-------------------------------------------------------------------
			# Step 3: For any duplicates found, delete the oldest one...
			#-------------------------------------------------------------------
			len=${#fmatch[@]}
			if [ ${len} -gt 0 ]; then
		        for (( i=0; i<${len}; i++ )); do
	                f1="${fmatch[$i]}.${ematch1[$i]}"
	                f2="${fmatch[$i]}.${ematch2[$i]}"
	                if [[ ${f1} -nt ${f2} ]]; then
                        echo "Keeping ${f1}, removing ${f2}"
                        rm ${f2}
	                else
                        echo "Keeping ${f2}, removing ${f1}"
                        rm ${f1}
	                fi      
		        done
			fi
			#===================
			#  Back to ${INST}
			#===================
			popd
		else
			echo "DOWNLOAD not needed for instance ${inst}"
		fi
		#======================
		#  Back to ${WORKDIR}
		#======================
		popd
	fi
done
#=========================
#  Back to phonebook dir
#=========================
popd

