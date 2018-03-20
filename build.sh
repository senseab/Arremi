#!/bin/sh

export BUILD_BASE=build

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
    cp assets/darwin/Info.plist $INSTALL_TARGET
}

function build_linux() {

}

function main() {
    mkdir -p $BUILD_BASE
    get_os
    case OS
        darwin)

        ;;
        linux)
        ;;
        *)
        echo "No target"
        ;;
    esac
}
