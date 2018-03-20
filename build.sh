#!/bin/sh

export BUILD_BASE=build

function get_go_ver(){
    GOVER=`go version|awk '{print $3}'`
    export GOVER
}

function get_os() {
    OS=`uname|tr '[:upper:]' '[:lower:]'`
    if [ "${GOOS}x" != "x" ]; then
        OS=$GOOS
    fi
    export OS
}

function build_darwin() {
    export INSTALL_TARGET=$BUILD_BASE/Arremi.app/Contents/
    mkdir -p $INSTALL_TARGET/MacOS
    go build -o $INSTALL_TARGET/MacOS/Arremi main.go
    mkdir -p $INSTALL_TARGET/Resources
    (cd assets/darwin/ && iconutil --convert icns Arremi.iconset --file ../../$INSTALL_TARGET/Resources/Arremi.icns)
    cp assets/darwin/Info.plist $INSTALL_TARGET
}

function build_linux() {

}

function main() {
    rm -rf build
    mkdir -p $BUILD_BASE
    get_os
    get_go_ver
    if [ $GOVER == 'go1.10' ]; then
        export CGO_LDFLAGS_ALLOW=".*"
    fi
    case $OS
        darwin)
            build_darwin
        ;;
        linux)
        ;;
        *)
        echo "No target"
        ;;
    esac
}
