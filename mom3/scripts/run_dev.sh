#!/bin/bash
echo "📦 Compilando..."
cd build || exit
cmake ..
make
echo "🚀 Ejecutando:"
./mom_core_v2
