#!/bin/bash

set -e

echo "=== Запуск Bash-пайплайна картирования ==="

REF="ncbi_dataset/data/GCF_000005845.2/GCF_000005845.2_ASM584v2_genomic.fna"
R1="SRR39029116_1.fastq"
R2="SRR39029116_2.fastq"

SAM="sample.sam"
BAM="sample.bam"
FLAGSTAT="flagstat.txt"
SORTED_BAM="sample.sorted.bam"
VCF="sample.vcf"

echo "Шаг 1: Картирование ридов с помощью BWA..."
bwa mem "$REF" "$R1" "$R2" > "$SAM"

echo "Шаг 2: Конвертация SAM в BAM..."
samtools view -bS "$SAM" > "$BAM"

echo "Шаг 3: Сбор статистики картирования..."
samtools flagstat "$BAM" > "$FLAGSTAT"

echo "Шаг 4: Проверка качества картирования..."

set +e
./parse_flagstat.sh "$FLAGSTAT" "$BAM" "$SORTED_BAM"
VERDICT_CODE=$?
set -e

if [ $VERDICT_CODE -eq 0 ]; then
    echo "Статус: OK. Продолжаем анализ вариантов (Variant Calling)..."

    echo "Шаг 5: Сортировка BAM-файла..."
    samtools sort "$SORTED_BAM" -o "$SORTED_BAM"
    
    echo "Шаг 6: Поиск генетических вариантов с помощью FreeBayes..."
    freebayes -f "$REF" "$SORTED_BAM" > "$VCF"
    
    echo "=== Пайплайн успешно завершен! Создан файл мутаций: $VCF ==="
else
    echo "Статус: not OK. Качество картирования ниже 90%. Пайплайн остановлен."
    exit 1
fi