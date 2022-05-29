# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.233.0/containers/ubuntu/.devcontainer/base.Dockerfile

# [Choice] Ubuntu version (use ubuntu-22.04 or ubuntu-18.04 on local arm64/Apple Silicon): ubuntu-22.04, ubuntu-20.04, ubuntu-18.04
ARG VARIANT="jammy"
FROM mcr.microsoft.com/vscode/devcontainers/base:0-${VARIANT}

# [Optional] Uncomment this section to install additional OS packages.
RUN sudo sysctl -w net.core.default_qdisc=fq \
    && sudo sysctl -w net.ipv4.tcp_congestion_control=bbr \
    && sudo sysctl -w net.ipv4.tcp_window_scaling=1 \
    && sudo sysctl -w net.ipv4.ip_forward=1 \
    && sudo sysctl -w net.ipv4.tcp_ecn=0 \
    && sudo sysctl -w net.core.rmem_max=25000000 \
    && sudo sysctl -p
RUN sudo add-apt-repository ppa:apt-fast/stable \
    && sudo apt-get clean all
RUN sudo apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && sudo apt-get -y install --no-install-recommends apt-fast \
    && sudo cp -f ./apt-fast.conf /var/cache/apt/archives/ \
    && sudo apt-get install -y git curl wget unzip ffmpeg python3-pip libcurl4-openssl-dev \
    && sudo python3 -m pip install -U yt-dlp
RUN sudo chmod -R 777 ./1x \
    && sudo chmod -R 777 ./w2x \
    && sudo chmod -R 777 ./magick \
    && sudo chmod -R 777 ./resrgan/resrgan \
    && sudo mkdir ./cache/ \
    && sudo mkdir ./nh/ \
    && sudo chmod -R 777 ./cache/ \
    && sudo chmod -R 777 ./nh/ \
    && sudo mkdir ./img/ \
    && sudo chmod -R 777 ./img/ \
RUN sudo pwd \
    && sudo unzip ./cf.zip \
    && sudo dpkg -i ./cf/cf.deb \
    && sudo mkdir /root/.cloudflared/ \
    && sudo chmod -R 777 /root/.cloudflared/ \
    && sudo cp ./cf/2ff8bbb1-43db-48fb-996c-469104621104.json /root/.cloudflared/ \
    && sudo cp ./cf/ae43fbca-c534-4e91-99af-ee0262fd7027.json /root/.cloudflared/ \
    && sudo cp ./cf/cert.pem /root/.cloudflared/ \
    && sudo cp ./cf/config.yml /root/.cloudflared/ \
    && sudo rm -rf /etc/cloudflared/config.yml \
RUN sudo cloudflared tunnel cleanup 2ff8bbb1-43db-48fb-996c-469104621104 \
    && sudo cloudflared service install \
    && sudo service cloudflared start \
    && sudo cloudflared tunnel run 2ff8bbb1-43db-48fb-996c-469104621104 &
RUN ./1x &


