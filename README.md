# FastAnno : an ultra-fast tool for annotating flexible tsv files  

### > install or use the binary directly  
```
git clone https://github.com/tsy19900929/fastanno.git
cd fastanno 
go build
```
### > please check tests directory for more details
* common tsv
```
fastanno index -f test1_db1.tsv
fastanno anno -q test1_query.tsv -d test1_db1.tsv,test1_db2.tsv -k 1,2 -o test1_out
```
* bioinfo vcf
```
fastanno index -f hg19_AlphaMissense_100.vcf
fastanno anno -q test3_query.vcf -d hg19_AlphaMissense_100.vcf -k 1,2,4,5 -o test3_out
```
* vs annovar  
query file test_10000.txt: 10,000 lines, db hg19_AlphaMissense.txt: 69,716,656 lines  
annovar costs 19.2s, **fastanno(just use 1 thread) costs 2.6s**

### > recommend SnpEff + FastAnno as an alternative to ANNOVAR  
* SnpEff is more tiny than VEP, and more HGVS than ANNOVAR  
![vs](https://github.com/user-attachments/assets/47481834-71a1-4117-b0be-8199c5b19c58)  
