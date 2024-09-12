#Only recommend using for large db, must already be sorted, use tests/sort_db.sh
./fastanno index -f hg19_AlphaMissense_100.vcf

./fastanno anno -q test3_query.vcf -d hg19_AlphaMissense_100.vcf -k 1,2,4,5 -o test3_out
