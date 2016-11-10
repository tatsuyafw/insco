#!/bin/bash

set -u

BINARY_DIR="$HOME/bin"
PREFIX_DIR="$HOME/usr/local"
SUPPORTS=("emacs" "ghq" "git" "peco" "vim")
TARGET=${1:-}
VERSION=${2:-}

FETCH_CMD=""

if type wget 2> /dev/null 1> /dev/null
then
    FETCH_CMD="wget -O"
elif type curl 2> /dev/null 1> /dev/null
then
    FETCH_CMD="curl -L -o"
else
    echo "[Error]: Install wget or curl."
    exit 1
fi

function usage() {
    cat <<EOS
Usage:
 $ $0 target [version]
EOS
}

function setup() {
    if [ ! -d $BINARY_DIR ]
    then
        mkdir $BINARY_DIR
    fi

    if [ ! -d $PREFIX_DIR ]
    then
        mkdir -p $PREFIX_DIR
    fi
}

function emacs() {
    TARGET_VERSION=${VERSION:-24.5}
    CONTENT="emacs-${TARGET_VERSION}"
    ARCH_FILE="${CONTENT}.tar.gz"
    MIRROR_LIST_URL="http://ftpmirror.gnu.org/emacs"
    FLAGS="--without-x"

    setup

    TMP_DIR=$(mktemp -d)
    ${FETCH_CMD} ${TMP_DIR}/${ARCH_FILE} ${MIRROR_LIST_URL}/${ARCH_FILE}
    tar zxf $TMP_DIR/$ARCH_FILE -C $TMP_DIR
    cd $TMP_DIR/$CONTENT

    ./configure --prefix=$PREFIX_DIR/$CONTENT $FLAGS
    make -j2
    make install

    if [ -f $BINARY_DIR/emacs ]
    then
        mv $BINARY_DIR/emacs $BINARY_DIR/emacs.org
    fi
    ln -s $PREFIX_DIR/$CONTENT/bin/emacs $BINARY_DIR/emacs
}

function ghq() {
    TARGET_VERSION=${VERSION:-0.7.4}
    CONTENT="ghq_linux_amd64"
    ARCH_FILE="${CONTENT}.zip"
    ARCH_URL="https://github.com/motemen/ghq/releases/download/v${TARGET_VERSION}"

    setup

    TMP_DIR=$(mktemp -d)
    ${FETCH_CMD} ${TMP_DIR}/${ARCH_FILE} ${ARCH_URL}/${ARCH_FILE}
    unzip $TMP_DIR/$ARCH_FILE -d $TMP_DIR

    mkdir $PREFIX_DIR/ghq-${TARGET_VERSION}
    cp $TMP_DIR/ghq $PREFIX_DIR/ghq-${TARGET_VERSION}

    if [ -f $BINARY_DIR/ghq -o -h $BINARY_DIR/ghq ]
    then
        mv $BINARY_DIR/ghq $BINARY_DIR/ghq.org
    fi
    ln -s $PREFIX_DIR/ghq-${TARGET_VERSION}/ghq $BINARY_DIR/ghq
}

function git() {
    TARGET_VERSION=${VERSION:-2.8.2}
    CONTENT="git-${TARGET_VERSION}"
    ARCH_FILE="v${TARGET_VERSION}.tar.gz"
    ARCH_URL="https://github.com/git/git/archive"
    FLAGS=""

    setup

    TMP_DIR=$(mktemp -d)
    ${FETCH_CMD} ${TMP_DIR}/${ARCH_FILE} ${ARCH_URL}/${ARCH_FILE}
    tar xf $TMP_DIR/$ARCH_FILE -C $TMP_DIR
    cd $TMP_DIR/$CONTENT

    make configure
    ./configure --prefix=$PREFIX_DIR/$CONTENT $FLAGS
    make -j2
    make install

    if [ -f $BINARY_DIR/git -o -h $BINARY_DIR/git ]
    then
        mv $BINARY_DIR/git $BINARY_DIR/git.org
    fi
    ln -s $PREFIX_DIR/$CONTENT/bin/git $BINARY_DIR/git
}

function peco() {
    TARGET_VERSION=${VERSION:-0.3.5}
    CONTENT="peco_linux_amd64"
    ARCH_FILE="${CONTENT}.tar.gz"
    ARCH_URL="https://github.com/peco/peco/releases/download/v${TARGET_VERSION}"

    setup

    TMP_DIR=$(mktemp -d)
    ${FETCH_CMD} ${TMP_DIR}/${ARCH_FILE} ${ARCH_URL}/${ARCH_FILE}
    tar xf $TMP_DIR/$ARCH_FILE -C $TMP_DIR
    cd $TMP_DIR/$CONTENT

    cp -r $TMP_DIR/$CONTENT $PREFIX_DIR/peco-${TARGET_VERSION}

    if [ -f $BINARY_DIR/peco -o -h $BINARY_DIR/peco ]
    then
        mv $BINARY_DIR/peco $BINARY_DIR/peco.org
    fi
    ln -s $PREFIX_DIR/peco-${TARGET_VERSION}/peco $BINARY_DIR/peco
}

function vim() {
    TARGET_VERSION=${VERSION:-7.4}
    CONTENT="vim-${TARGET_VERSION}"
    ARCH_FILE="${CONTENT}.tar.bz2"
    ARCH_DIR="vim"$(echo $TARGET_VERSION | tr -d '.')
    ARCH_URL="http://ftp.vim.org/pub/vim/unix"
    FLAGS=""

    setup

    TMP_DIR=$(mktemp -d)
    ${FETCH_CMD} ${TMP_DIR}/${ARCH_FILE} ${ARCH_URL}/${ARCH_FILE}
    tar xf $TMP_DIR/$ARCH_FILE -C $TMP_DIR
    cd $TMP_DIR/$ARCH_DIR

    ./configure --prefix=$PREFIX_DIR/$ARCH_DIR $FLAGS
    make -j2
    make install

    if [ -f $BINARY_DIR/vim -o -h $BINARY_DIR/vim ]
    then
        mv $BINARY_DIR/vim $BINARY_DIR/vim.org
    fi
    ln -s $PREFIX_DIR/$ARCH_DIR/bin/vim $BINARY_DIR/vim
}

if [ "$TARGET" == "" ]
then
    usage
    exit 0
fi

case $TARGET in
    emacs)
        emacs
        ;;
    ghq)
        ghq
        ;;
    git)
        git
        ;;
    peco)
        peco
        ;;
    vim)
        vim
        ;;
    *)
        echo "[Error] Not support $TARGET"
        echo "  Support: "
        for s in "${SUPPORTS[@]}"
        do
            echo "    $s"
        done
        exit 1
        ;;
esac
