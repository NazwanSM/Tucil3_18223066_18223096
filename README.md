# Tucil3_18223066_18223096 — Ice Sliding Puzzle Solver

Penyelesai *Ice Sliding Puzzle* menggunakan algoritma pencarian
informed dan uninformed: **Uniform Cost Search (UCS)**, **Greedy
Best-First Search (GBFS)**, dan **A\***.

## Aturan Singkat

- Aktor meluncur ke atas / bawah / kiri / kanan dan baru berhenti
  tepat sebelum menabrak batu (`X`).
- Aktor harus menginjak petak angka (`0`..`9`) secara **berurutan**
  sebelum sampai di petak keluar (`O`).
- Meluncur menembus tepi papan tanpa terhalang batu = *game over*.
- Menyentuh lava (`L`) = *game over*.
- *Cost* satu luncuran = jumlah cost setiap petak yang dilewati,
  **kecuali** petak awal luncuran. Petak berhenti tetap dihitung.

## Struktur Direktori

```
Tucil3_18223066_18223096/
├── README.md
├── go.mod
├── bin/                 # binary hasil build (di-ignore git)
├── doc/                 # laporan PDF
├── test/                # kumpulan kasus uji *.txt
└── src/
    ├── main.go          # entry point CLI
    ├── board/           # representasi papan + parser
    ├── state/           # State (posisi + target digit)
    ├── successor/       # simulator luncuran -> generator suksesor
    ├── search/          # core best-first search + Node + PQ
    └── heuristic/       # h(n): Zero, Manhattan, RemainingDigits, ...
```

## Build & Run

```powershell
# build ke folder bin/
go build -o bin/iceslide.exe ./src

# run
./bin/iceslide.exe <input.txt>
```

