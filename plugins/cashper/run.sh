# ! /user/bin/bash
CURDIR=`pwd`
CURDIR=`basename ${CURDIR}`
echo "CURR_DIR: $CURDIR"
GOFILE=`ls ${CURDIR}*.go`
echo "GOFILE: $GOFILE"
FILENAME="${GOFILE%.*}"
echo "FILENAME: $FILENAME"

go build -buildmode=plugin -o ${FILENAME}.so ${GOFILE}
