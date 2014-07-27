
sudo aptitude install -y nginx python python-pip python-dev golang-go
sudo pip install virtualenv

if [ ! -s `which go` ]; then
  wget http://golang.org/dl/go1.3.linux-amd64.tar.gz
  tar -C /usr/local -xzf go1.3.linux-amd64.tar.gz
fi

WORKDIR=/vagrant

if [ ! -d $WORKDIR/venv ]; then
  virtualenv $WORKDIR/venv
fi

cd $WORKDIR
source ./venv/bin/activate
pip install fabric
