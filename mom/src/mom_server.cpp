#include <iostream>
#include <memory>
#include <string>
#include <grpcpp/grpcpp.h>
#include "mom.grpc.pb.h"
#include "broker.hpp"

#include <ifaddrs.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <unistd.h>

using grpc::Server;
using grpc::ServerBuilder;
using grpc::ServerContext;
using grpc::Status;

class MomServiceImpl final : public mom::MomService::Service {
private:
    Broker broker;

public:

    MomServiceImpl(int port) : broker("/home/carlosm/Documentos/Projects/Proyecto1_TopicosTelematica/mom/data/mom" + std::to_string(port) + ".db") {}

    Status RegistrarUsuario(ServerContext*, const mom::Credenciales* req, mom::RespuestaSimple* res) override {
        bool ok = broker.registrarUsuario(req->username(), req->password());
        res->set_exito(ok);
        res->set_mensaje(ok ? "Usuario registrado" : "Error al registrar");
        return Status::OK;
    }

    Status AutenticarUsuario(ServerContext*, const mom::Credenciales* req, mom::Token* res) override {
        std::string token;
        bool ok = broker.autenticarUsuario(req->username(), req->password(), token);
        if (ok) res->set_token(token);
        return Status::OK;
    }

    Status GuardarTokenReplica(ServerContext*, const mom::TokenConExpiracion* req, mom::RespuestaSimple* res) override {
        bool ok = broker.guardarTokenReplica(req->username(), req->token(), req->expiracion());
        res->set_exito(ok);
        res->set_mensaje(ok ? "Token replicado" : "Error");
        return Status::OK;
    }

    Status CrearCola(ServerContext*, const mom::AccionConToken* req, mom::RespuestaSimple* res) override {
        bool ok = broker.crearCola(req->nombre(), req->token());
        res->set_exito(ok);
        res->set_mensaje(ok ? "Cola creada" : "Error");
        return Status::OK;
    }

    Status EliminarCola(ServerContext*, const mom::AccionConToken* req, mom::RespuestaSimple* res) override {
        bool ok = broker.eliminarCola(req->nombre(), req->token());
        res->set_exito(ok);
        res->set_mensaje(ok ? "Cola eliminada" : "Error");
        return Status::OK;
    }

    Status AutorizarUsuarioCola(ServerContext*, const mom::AutorizacionColaRequest* req, mom::RespuestaSimple* res) override {
        bool ok = broker.autorizarCola(req->nombre(), req->usuarioobjetivo(), req->token());
        res->set_exito(ok);
        res->set_mensaje(ok ? "Usuario autorizado" : "Error");
        return Status::OK;
    }

    Status EnviarMensajeCola(ServerContext*, const mom::MensajeConToken* req, mom::RespuestaSimple* res) override {
        bool ok = broker.enviarMensajeACola(req->nombre(), req->contenido(), req->token());
        res->set_exito(ok);
        res->set_mensaje(ok ? "Mensaje enviado" : "Error");
        return Status::OK;
    }

    Status ConsumirMensajeCola(ServerContext*, const mom::AccionConToken* req, mom::MensajeTexto* res) override {
        auto msg = broker.consumirMensajeDeCola(req->nombre(), req->token());
        if (msg.has_value()) {
            res->set_contenido(msg->getContenido());
            res->set_remitente(msg->getRemitente());
            res->set_canal(msg->getCanal());
            res->set_timestamp("...");
        }
        return Status::OK;
    }

    Status CrearTopico(ServerContext*, const mom::AccionConToken* req, mom::RespuestaSimple* res) override {
        bool ok = broker.crearTopico(req->nombre(), req->token());
        res->set_exito(ok);
        res->set_mensaje(ok ? "Tópico creado" : "Error");
        return Status::OK;
    }

    Status EliminarTopico(ServerContext*, const mom::AccionConToken* req, mom::RespuestaSimple* res) override {
        bool ok = broker.eliminarTopico(req->nombre(), req->token());
        res->set_exito(ok);
        res->set_mensaje(ok ? "Tópico eliminado" : "Error");
        return Status::OK;
    }

    Status SuscribirseTopico(ServerContext*, const mom::AccionConToken* req, mom::RespuestaSimple* res) override {
        bool ok = broker.suscribirATopico(req->nombre(), req->token());
        res->set_exito(ok);
        res->set_mensaje(ok ? "Suscripción exitosa" : "Error");
        return Status::OK;
    }

    Status PublicarEnTopico(ServerContext*, const mom::MensajeConToken* req, mom::RespuestaSimple* res) override {
        bool ok = broker.publicarEnTopico(req->nombre(), req->contenido(), req->token());
        res->set_exito(ok);
        res->set_mensaje(ok ? "Mensaje publicado" : "Error");
        return Status::OK;
    }

    Status ConsumirDesdeTopico(ServerContext*, const mom::AccionConToken* req, mom::ListaMensajes* res) override {
        auto mensajes = broker.consumirDesdeTopico(req->nombre(), req->token());
        for (const auto& m : mensajes) {
            mom::MensajeTexto* nuevo = res->add_mensajes();
            nuevo->set_contenido(m.getContenido());
            nuevo->set_remitente(m.getRemitente());
            nuevo->set_canal(m.getCanal());
            nuevo->set_timestamp("...");
        }
        return Status::OK;
    }

    Status ListarColas(ServerContext*, const mom::Token* req, mom::ListaNombres* res) override {
        auto nombres = broker.listarColas();
        for (const auto& nombre : nombres) {
            res->add_nombres(nombre);
        }
        return Status::OK;
    }
    
    Status ListarTopicos(ServerContext*, const mom::Token* req, mom::ListaNombres* res) override {
        auto nombres = broker.listarTopicos();
        for (const auto& nombre : nombres) {
            res->add_nombres(nombre);
        }
        return Status::OK;
    }
};

void RunServer(int port) {
    std::string local_ip = "0.0.0.0"; 
    struct ifaddrs* ifaddr;
    if (getifaddrs(&ifaddr) == 0) {
        for (struct ifaddrs* ifa = ifaddr; ifa != nullptr; ifa = ifa->ifa_next) {
            if (ifa->ifa_addr && ifa->ifa_addr->sa_family == AF_INET) {
                char ip[INET_ADDRSTRLEN];
                void* addr = &((struct sockaddr_in*)ifa->ifa_addr)->sin_addr;
                inet_ntop(AF_INET, addr, ip, INET_ADDRSTRLEN);
                if (std::string(ifa->ifa_name) != "lo") { 
                    local_ip = ip;
                    break;
                }
            }
        }
        freeifaddrs(ifaddr);
    }

    std::string server_address(local_ip + ":" + std::to_string(port));
    MomServiceImpl service(port);

    grpc::ServerBuilder builder;
    builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());
    builder.RegisterService(&service);
    std::unique_ptr<grpc::Server> server(builder.BuildAndStart());
    std::cout << "Servidor MOM escuchando en " << server_address << std::endl;
    server->Wait();
}

bool verificarPuerto(int port) {
    int sock = socket(AF_INET, SOCK_STREAM, 0);
    if (sock < 0) return true;

    sockaddr_in addr{};
    addr.sin_family = AF_INET;
    addr.sin_addr.s_addr = INADDR_ANY; 
    addr.sin_port = htons(port);

    bool inUse = bind(sock, (struct sockaddr*)&addr, sizeof(addr)) < 0;
    close(sock);
    return inUse;
}

int main(int argc, char** argv) {
    int port = 50051;

    while (verificarPuerto(port)) {
        port = port + 1;

    }
    RunServer(port);
    return 0;
}