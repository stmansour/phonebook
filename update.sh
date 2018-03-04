#!/bin/bash
#############################################################################
# readConfig
#   Description:
#       Read the config.json file from the directory containing this script
#       to set some of the key values needed to access the artifactory repo.
#
#       Upon returning, URLBASE will end with the character "/".
#
#   Params:
#       none
#
#   Returns:
#       sets variables:  APIKEY, REPOUSER, and URLBASE
#
#############################################################################
readConfig() {
    RELDIR=$(cd `dirname "${BASH_SOURCE[0]}"` && pwd)
    CONF="${RELDIR}/config.json"
    REPOUSER=$(grep RepoUser ${CONF} | awk '{print $2;}' | sed -e 's/[,"]//g')
    APIKEY=$(grep RepoPass ${CONF} | awk '{print $2;}' | sed -e 's/[,"]//g')
    URLBASE=$(grep RepoURL ${CONF} | awk '{print $2;}' | sed -e 's/[,"]//g')

    # add a trailing / if it does not have one...
    if [ "${URLBASE: -1}" != "/" ]; then
        URLBASE="${URLBASE}/"
    fi

    pushd
    cd ~ec2-user
    EC2USERHOME=$(pwd)
    popd

    echo "RELDIR = ${RELDIR}"
    echo "REPOUSER = ${REPOUSER}"
    echo "APIKEY = ${APIKEY}"
    echo "URLBASE = ${URLBASE}"
    echo "EC2USERHOME = ${EC2USERHOME}"
}

#############################################################################
# configure
#   Description:
#       The only configuration needed is the jfrog cli environment. Just
#       make sure we have it in the path. If it is not present, then
#       get it.
#
#   Params:
#       none
#
#   Returns:
#       nothing
#
#############################################################################
configure() {
    #---------------------------------------------------------
    # the user's bin directory is not created by default...
    #---------------------------------------------------------
    if [ ! -d ${EC2USERHOME}/bin ]; then
        mkdir ${EC2USERHOME}/bin
    fi

    #---------------------------------------------------------
    # now make sure that we have jfrog...
    #---------------------------------------------------------
    if [ ! -f ${EC2USERHOME}/bin/jfrog ]; then
        curl -s -u "${REPOUSER}:${APIKEY}" ${URLBASE}accord/tools/jfrog > ${EC2USERHOME}/bin/jfrog
        chown ec2-user:ec2-user ${EC2USERHOME}/bin/jfrog
        chmod +x ${EC2USERHOME}/bin/jfrog
    fi
    if [ ! -d ${EC2USERHOME}/.jfrog ]; then
        curl -s -u "${REPOUSER}:${APIKEY}" ${URLBASE}accord/tools/jfrogconf.tar > ${EC2USERHOME}/jfrogconf.tar
        pushd ${EC2USERHOME}
        tar xvf jfrogconf.tar
        rm jfrogconf.tar
        chown ec2-user:ec2-user ${EC2USERHOME}/bin/jfrog
        popd
    fi
    if [ ! -d ~root/.jfrog ]; then
        curl -s -u "${REPOUSER}:${APIKEY}" ${URLBASE}accord/tools/jfrogconf.tar > ~root/jfrogconf.tar
        pushd ~root
        tar xvf jfrogconf.tar
        rm jfrogconf.tar
        popd
    fi
    JFROG="${EC2USERHOME}/bin/jfrog"
    echo "JFROG = ${JFROG}"
}


#############################################################################
# GetLatestProductRelease
#   Description:
#       The only configuration needed is the jfrog cli environment. Just
#       make sure we have it in the path. If it is not present, then
#       get it.
#
#   Params:
#       ${1} = base name of product (rentroll, phonebook, mojo, ...)
#
#   Returns:
#       nothing
#
#############################################################################
GetLatestRepoRelease() {
    echo "GetLatestRepoRelease:  looking for released ${1}"
    f=$(${JFROG} rt s "accord/air/release/*" | grep ${1} | awk '{print $2}' | sed 's/"//g')
    echo "f = ${f}"
    if [ "x${f}" = "x" ]; then
        echo "There are no product releases for ${f}"
        exit 1
    fi
    t=$(basename ${f})
    echo "t = ${t}"
    cd ${RELDIR}
    cdir=$(pwd)
    echo "calling curl in directory ${cdir}"
    echo "curl -s -u ${REPOUSER}:${APIKEY} ${URLBASE}${f} > ../${t}"
    curl -s -u "${REPOUSER}:${APIKEY}" ${URLBASE}${f} > ../${t}
    echo "After call to curl, directory contents ls .."
    ls ..
}


readConfig
configure

#----------------------------------------------
#  ensure that we're in the phonebook directory...
#----------------------------------------------

cd ${RELDIR}
dir=${PWD##*/}
if [ ${dir} != "phonebook" ]; then
    echo "This script must execute in the phonebook directory."
    echo "current directory is: ${dir}"
    exit 1
fi

user=$(whoami)
if [ ${user} != "root" ]; then
    echo "This script must execute as root.  Try sudo !!"
    exit 1
fi

echo -n "Shutting down phonebook server."; $(./activate.sh stop) >/dev/null 2>&1
echo -n "."
echo -n "."; sleep 6
echo -n "."; cd ..
echo -n "."; rm -f phonbook*.tar*
echo " "

echo -n "Retrieving latest released Phonebook..."
GetLatestRepoRelease "phonebook"

echo "Installing.."
echo -n "."; cd ${RELDIR}/..
echo -n "."; gunzip -f phonebook*.tar.gz
echo -n "."; chown -R ec2-user:ec2-user rentroll
#echo -n "."; rm -f phonebook*.tar*
echo -n "."; cd ${RELDIR}
echo

# echo -n "."; updateImages
echo -n "."; chmod u+s phonebook pbwatchdog
echo -n "."; echo -n "starting..."
echo -n "."; ./activate.sh -b start
echo -n "."; sleep 2
echo -n "."; status=$(./activate.sh ready)
if [ "${status}" = "OK" ]; then
    echo "Activation successful"
    rm ../phonebook*.tar
else
    echo "Problems activating phonebook.  Status = ${status}"
fi
