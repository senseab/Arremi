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
    mkdir -p $INSTALL_TARGET/Resources
    (cd assets/darwin/ && iconutil --convert icns Arremi.iconset --file ../../$INSTALL_TARGET/Resources/Arremi.icns)
}

function build_linux() {

}

function main() {
    mkdir -p $BUILD_BASE
    get_os
    case OS
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
