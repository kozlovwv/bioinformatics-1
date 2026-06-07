package main

import (
	"os"
	"path/filepath"

	"github.com/scipipe/scipipe"
)

func main() {
	wf := scipipe.NewWorkflow("dna_mapping_pipeline", 4)

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	refPath := filepath.Join(dir, "ncbi_dataset/data/GCF_000005845.2/GCF_000005845.2_ASM584v2_genomic.fna")
	r1Path  := filepath.Join(dir, "SRR39029116_1.fastq")
	r2Path  := filepath.Join(dir, "SRR39029116_2.fastq")
	script  := filepath.Join(dir, "parse_flagstat.sh")

	bwaAlign := wf.NewProc("bwa_align", "bwa mem "+refPath+" "+r1Path+" "+r2Path+" > {o:sam}")
	bwaAlign.SetOut("sam", "sample.sam")

	samToBam := wf.NewProc("sam_to_bam", "samtools view -bS {i:sam} > {o:bam}")
	samToBam.SetOut("bam", "sample.bam")

	flagstat := wf.NewProc("samtools_flagstat", "samtools flagstat {i:bam} > {o:txt}")
	flagstat.SetOut("txt", "flagstat.txt")

	parser := wf.NewProc("parse_flagstat", script+" {i:txt} {i:bam} {o:passed_bam} > verdict.txt")
	parser.SetOut("passed_bam", "sample.passed.bam")

	samSort := wf.NewProc("samtools_sort", "samtools sort {i:passed_bam} -o {o:sorted_bam}")
	samSort.SetOut("sorted_bam", "sample.sorted.bam")

	freebayes := wf.NewProc("freebayes", "freebayes -f "+refPath+" {i:sorted_bam} > {o:vcf}")
	freebayes.SetOut("vcf", "sample.vcf")

	samToBam.In("sam").From(bwaAlign.Out("sam"))
	flagstat.In("bam").From(samToBam.Out("bam"))
	
	parser.In("txt").From(flagstat.Out("txt"))
	parser.In("bam").From(samToBam.Out("bam"))

	samSort.In("passed_bam").From(parser.Out("passed_bam"))
	freebayes.In("sorted_bam").From(samSort.Out("sorted_bam"))

	wf.PlotGraph("pipeline_dag.dot")
	wf.Run()
}