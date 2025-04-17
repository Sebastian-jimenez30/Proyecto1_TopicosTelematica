cmake_minimum_required(VERSION 3.10)
project(MOMCoreV2)

set(CMAKE_CXX_STANDARD 17)

# -------------------- DEPENDENCIAS --------------------
find_package(OpenSSL REQUIRED)
find_package(SQLite3 REQUIRED)
find_package(Protobuf REQUIRED)
find_package(gRPC REQUIRED)

# -------------------- VERIFICAR PLUGIN DE gRPC --------------------
find_program(GRPC_CPP_PLUGIN grpc_cpp_plugin PATHS /usr/local/bin /usr/bin)
if (NOT GRPC_CPP_PLUGIN)
    message(FATAL_ERROR "gRPC plugin not found. Please ensure grpc_cpp_plugin is installed.")
endif()

# -------------------- INCLUDES --------------------
include_directories(
    ${PROJECT_SOURCE_DIR}/include
    ${PROJECT_SOURCE_DIR}/external/jwt-cpp/include
    ${OPENSSL_INCLUDE_DIR}
    ${SQLite3_INCLUDE_DIRS}
    ${Protobuf_INCLUDE_DIRS}
    ${gRPC_INCLUDE_DIRS}
    ${PROJECT_SOURCE_DIR}/src
)

# -------------------- GENERAR CÓDIGO gRPC Y PROTOBUF --------------------
file(GLOB proto_files "${PROJECT_SOURCE_DIR}/proto/*.proto")

add_custom_command(
  OUTPUT
    ${PROJECT_SOURCE_DIR}/src/mom.pb.cc
    ${PROJECT_SOURCE_DIR}/src/mom.pb.h
    ${PROJECT_SOURCE_DIR}/src/mom.grpc.pb.cc
    ${PROJECT_SOURCE_DIR}/src/mom.grpc.pb.h
  COMMAND protoc
    --proto_path=${PROJECT_SOURCE_DIR}/proto
    --cpp_out=${PROJECT_SOURCE_DIR}/src
    --grpc_out=${PROJECT_SOURCE_DIR}/src
    --plugin=protoc-gen-grpc=${GRPC_CPP_PLUGIN}
    ${proto_files}
  DEPENDS ${proto_files}
)

add_custom_target(proto_files_target DEPENDS
  ${PROJECT_SOURCE_DIR}/src/mom.pb.cc
  ${PROJECT_SOURCE_DIR}/src/mom.pb.h
  ${PROJECT_SOURCE_DIR}/src/mom.grpc.pb.cc
  ${PROJECT_SOURCE_DIR}/src/mom.grpc.pb.h
)

# -------------------- EXECUTABLE: grpc_server --------------------
add_executable(grpc_server
  src/mom_server.cpp
  src/mom.pb.cc
  src/mom.grpc.pb.cc
  src/broker.cpp
  src/usuario.cpp
  src/mensaje.cpp
  src/persistencia.cpp
  src/cola.cpp
  src/topico.cpp
)

add_dependencies(grpc_server proto_files_target)

target_link_libraries(grpc_server
  gRPC::grpc++
  protobuf::libprotobuf
  OpenSSL::SSL
  OpenSSL::Crypto
  SQLite::SQLite3
)
