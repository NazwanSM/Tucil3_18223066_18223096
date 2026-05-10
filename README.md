# Ice Sliding Puzzle Solver

Tucil3 IF2211 Strategi Algoritma — Semester II 2025/2026

## Deskripsi Program

Program penyelesai *Ice Sliding Puzzle*: sebuah permainan logika di mana aktor meluncur di atas es dan hanya berhenti saat menabrak batu atau dinding. Tujuan aktor adalah mencapai petak keluar (`O`) setelah melewati semua petak angka (`0`–`9`) secara berurutan.

Program mengimplementasikan empat algoritma *pathfinding*:

| Algoritma | Keterangan |
|-----------|-----------|
| **UCS** | Uniform Cost Search — uninformed, menjamin solusi optimal |
| **GBFS** | Greedy Best-First Search — informed, cepat tapi tidak selalu optimal |
| **A\*** | A-Star — informed, optimal dengan heuristik admissible |
| **IDA\*** | Iterative Deepening A-Star — hemat memori, optimal (bonus) |

Lima pilihan heuristik tersedia untuk GBFS, A\*, dan IDA\*:

- **Manhattan** — jarak Manhattan ke target berikutnya
- **Remaining Digits** — jumlah digit yang belum dikumpulkan
- **Euclidean** — jarak Euclidean ke target berikutnya (bonus)
- **Max Combined** — max(Manhattan, Remaining Digits) (bonus)
- **Zero** — tanpa heuristik, dipakai otomatis oleh UCS

Program tersedia dalam dua mode: **CLI** (terminal interaktif) dan **GUI** (antarmuka web di browser).

## Requirement

- **Go** versi 1.22 atau lebih baru
- Tidak memerlukan instalasi library tambahan (tidak ada CGO)
- Untuk mode GUI: browser web (Chrome, Firefox, Edge, dll.)

Cek versi Go:

```bash
go version
```

Unduh Go: https://go.dev/dl/

## Cara Kompilasi

Clone repository lalu build:

```bash
# Windows
go build -o bin/iceslide.exe ./src

# Linux / macOS
go build -o bin/iceslide ./src
```

Atau gunakan flag `CGO_ENABLED=0` untuk memastikan build tanpa dependensi C:

```bash
CGO_ENABLED=0 go build -o bin/iceslide.exe ./src
```

Binary hasil build tersimpan di folder `bin/`.

## Cara Menjalankan

### Mode GUI (disarankan)

```bash
# Windows
bin\iceslide.exe -gui

# Linux / macOS
./bin/iceslide -gui
```

Browser akan terbuka otomatis ke `http://127.0.0.1:8765`. Jika tidak terbuka, buka URL tersebut secara manual.

**Alur penggunaan GUI:**
1. Klik **Browse File (.txt)** atau drag & drop file puzzle ke area input
2. Pilih **Algoritma** (UCS / GBFS / A\* / IDA\*)
3. Pilih **Heuristik** (untuk GBFS / A\* / IDA\*)
4. Klik **Cari Solusi**
5. Setelah selesai, gunakan kontrol playback untuk melihat langkah demi langkah:
   - Tombol First / Prev / Play / Next / Last untuk navigasi antar step
   - Slider untuk loncat ke step tertentu
   - Slider kecepatan (0.25× – 4×) untuk mengatur kecepatan animasi
6. Klik **Simpan Solusi (.txt)** untuk mengunduh solusi ke file

### Mode CLI

```bash
# Jalankan interaktif (pilih algoritma & heuristik via terminal)
bin\iceslide.exe test\01_basic.txt

# Jalankan dengan flag langsung
bin\iceslide.exe -alg astar -h manhattan test\01_basic.txt

# Simpan solusi ke file tertentu
bin\iceslide.exe -alg ucs -out hasil.txt test\02_one_digit.txt
```

Flag yang tersedia:

| Flag | Nilai | Keterangan |
|------|-------|-----------|
| `-alg` | `ucs` / `gbfs` / `astar` / `idastar` | Algoritma (opsional, tanya interaktif jika kosong) |
| `-h` | `manhattan` / `remaining` / `euclidean` / `maxcombined` | Heuristik |
| `-out` | path file | Path default untuk simpan solusi (default: `test/solusi.txt`) |
| `-gui` | — | Jalankan mode GUI |

**Contoh output CLI:**

```
== Initial Board ==
XXXXXXX
X0****X
X**X**X
X****OX
X1***LX
XZ**X*X
XXXXXXX

== Solution Summary ==
Algorithm       : A*
Heuristic       : Manhattan
Path            : RULUDRUR
Cost            : 87
Steps           : 8
Iterations      : 42
Nodes generated : 65
Duration        : 1.234 ms
```

### Format File Input

```
N M
<baris papan 1>
...
<baris papan N>
<cost baris 1>
...
<cost baris N>
```

Keterangan simbol:

| Simbol | Keterangan |
|--------|-----------|
| `Z` | Posisi awal aktor |
| `*` | Petak es (dapat dilewati) |
| `X` | Batu / rintangan (aktor berhenti sebelumnya) |
| `L` | Lava (game over jika dilewati) |
| `O` | Petak tujuan |
| `0`–`9` | Petak angka (harus dilalui berurutan) |

## Struktur Direktori

```
Tucil3_18223066_18223096/
├── README.md
├── go.mod
├── bin/                  # binary hasil build
├── doc/                  # laporan PDF
├── test/                 # kasus uji (.txt)
└── src/
    ├── main.go           # entry point (CLI + GUI)
    ├── board/            # representasi papan & parser
    ├── state/            # State (posisi + next digit)
    ├── successor/        # simulator luncuran & generator suksesor
    ├── search/           # algoritma pencarian (UCS/GBFS/A*/IDA*)
    ├── heuristic/        # fungsi heuristik
    ├── visualizer/       # output CLI & playback
    └── gui/              # server GUI berbasis web
```

## Author

| Nama | NIM | Kontak |
|------|-----|--------|
| Nazwan Siddqi Muttaqin | 18223066 | 18223066@std.stei.itb.ac.id |
| Matthew Sebastian Kurniawan | 18223096 | 18223096@std.stei.itb.ac.id |
