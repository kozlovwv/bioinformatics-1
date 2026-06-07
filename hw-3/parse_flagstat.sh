#!/bin/bash

FLAGSTAT_FILE=$1
BAM_FILE=$2
OUTPUT_PASSED_BAM=$3

PERCENT=$(grep " mapped (" "$FLAGSTAT_FILE" | \
head -n1 | \
sed -E 's/.*\(([0-9.]+)%.*/\1/')

echo "Mapped: $PERCENT %"

if awk "BEGIN {exit !($PERCENT > 90)}"
then
    echo "OK"
    cp "$BAM_FILE" "$OUTPUT_PASSED_BAM"
else
    echo "not OK"
    exit 1
fi