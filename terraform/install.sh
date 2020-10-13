wget https://golang.org/dl/go1.15.2.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.15.2.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
rm -f go1.15.2.linux-amd64.tar.gz

sudo apt install make -y
git clone https://github.com/IonutCraciun/databaseRestApi.git
cd databaseRestApi && make servicemachine && make start

