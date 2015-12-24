PASS=AP3wHZhcQQCvkC4GVCCZzPcqe3L
ART=http://ec2-52-91-201-195.compute-1.amazonaws.com/artifactory
USR=accord

EXTERNAL_HOST_NAME=$( curl http://169.254.169.254/latest/meta-data/public-hostname )
${EXTERNAL_HOST_NAME:?"Need to set EXTERNAL_HOST_NAME non-empty"}

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
    pushd phonebook
    /usr/local/accord/bin/getfile.sh jenkins-snapshot/phonebook/latest/pbimages.tar.gz
    gunzip -f pbimages.tar.gz
    tar xvf pbimages.tar
    popd
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
#  Now download the requested apps...
#----------------------------------------------
# - - - - -  APPEND DATA and DOWNLOAD APPS  - - - - - - -
# install_mysql
# UHURA_MASTER_URL=http://ip-172-31-56-33:8251/
# MY_INSTANCE_NAME="phone"
# mkdir ~ec2-user/apps;cd apps
# mkdir ~ec2-user/apps/tgo
# mkdir ~ec2-user/apps/phonebook
# artf_get jenkins-snapshot/tgo/latest tgo.tar.gz
artf_get jenkins-snapshot/phonebook/latest phonebook.tar.gz
# gunzip tgo.tar.gz;tar xf tgo.tar
gunzip phonebook.tar.gz;tar xf phonebook.tar
updateImages

