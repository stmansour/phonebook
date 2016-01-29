PASS=AP3wHZhcQQCvkC4GVCCZzPcqe3L
ART=http://ec2-52-91-201-195.compute-1.amazonaws.com/artifactory
USR=accord

EXTERNAL_HOST_NAME=$(curl -s http://169.254.169.254/latest/meta-data/public-hostname)
#${EXTERNAL_HOST_NAME:?"Need to set EXTERNAL_HOST_NAME non-empty"}

#--------------------------------------------------------------
#  Routine to download files from Artifactory
#--------------------------------------------------------------
artf_get() {
    echo "Downloading $1/$2"
    wget -O "$2" --user=$USR --password=$PASS ${ART}/"$1"/"$2"
}

#--------------------------------------------------------------
#  function to install mysql
#--------------------------------------------------------------
install_mysql() {
        echo "installing mysql"
        yum -y install mysql55-server.x86_64
        service mysqld start
        echo "CREATE DATABASE accord;use accord;GRANT ALL PRIVILEGES ON accord.* TO 'ec2-user'@'localhost';"  | mysql
}

restoredb() {
    echo "IN RESTOREDB"
    pushd /tmp
    DIR=$(pwd)
    echo "CURRENT DIRECTORY = ${DIR}"
    echo "${ACCORDHOME}/bin/getfile.sh getfile.sh ext-tools/testing/$1"
    ${ACCORDHOME}/bin/getfile.sh getfile.sh ext-tools/testing/$1
    echo "Get file $1 completed"
    echo "${ACCORDHOME}/testtools/restoredb.sh /tmp/$1"
    ${ACCORDHOME}/testtools/restoredb.sh /tmp/$1
    echo "restoredb.sh completed"
    popd
    DIR=$(pwd)
    echo "popd completed, dir = ${DIR}"
}

updatePkgs() {
    #--------------------------------------------------------------
    #  update all the out-of-date packages, add Java 1.8, and md5sum
    #--------------------------------------------------------------
    yum -y update
    yum -y install java-1.8.0-openjdk-devel.x86_64
    yum -y install isomd5sum.x86_64
}

updateImages() {
    /usr/local/accord/bin/getfile.sh jenkins-snapshot/phonebook/latest/pbimages.tar.gz
    rm -rf images
    gunzip -f pbimages.tar.gz
    tar xvf pbimages.tar
}

loadAccordTools() {
    #--------------------------------------------------------------
    #  Let's get our tools in place...
    #--------------------------------------------------------------
    artf_get ext-tools/utils accord-linux.tar.gz
    echo "Installing /usr/local/accord" >>${LOGFILE}
    cd /usr/local
    tar xvzf ~ec2-user/accord-linux.tar.gz
    chown -R ec2-user:ec2-user accord
    cd ~ec2-user/
}

#----------------------------------------------
#  ensure that we're in the phonebook directory...
#----------------------------------------------

dir=${PWD##*/}
if [ ${dir} != "phonebook" ]; then
    echo "This script must execute in the phonebook directory."
    exit 1
fi

user=$(whoami)
if [ ${user} != "root" ]; then
    echo "This script must execute as root.  Try sudo !!"
    exit 1
fi

$(./activate.sh stop)
echo "shutdown initiated..."
sleep 6
cd ..
echo "Retrieving latest phonebook..."
/usr/local/accord/bin/getfile.sh jenkins-snapshot/phonebook/latest/phonebook.tar.gz
# gunzip tgo.tar.gz;tar xf tgo.tar
gunzip -f phonebook.tar.gz;tar xvf phonebook.tar
chown -R ec2-user:ec2-user phonebook
cd phonebook/
updateImages
chmod u+s phonebook pbwatchdog
echo "starting..."
./activate.sh -b start
sleep 3
status=$(./activate.sh ready)
if [ ${status} == "OK" ]; then
    echo "Activation successful"
else
    echo "Problems activating phonebook.  Status = ${status}"
fi