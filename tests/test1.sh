#Only recommend using for large db, must already be sorted, use tests/sort_db.sh
./fastanno index -f test1_db1.tsv

./fastanno anno -q test1_query.tsv -d test1_db1.tsv,test1_db2.tsv -k 1,2 -o test1_out
