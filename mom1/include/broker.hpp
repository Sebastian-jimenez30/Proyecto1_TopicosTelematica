#ifndef BROKER_HPP
#define BROKER_HPP

#include <unordered_map>
#include <optional>
#include <memory>
#include "usuario.hpp"
#include "cola.hpp"
#include "topico.hpp"
#include "persistencia.hpp"

class Broker {
public:
    Broker(const std::string& db_path);

    // --- Usuarios
    bool registrarUsuario(const std::string& username, const std::string& password);
    bool autenticarUsuario(const std::string& username, const std::string& password, std::string& token);
    bool verificarToken(const std::string& token, std::string& username);

    // --- Colas
    bool crearCola(const std::string& nombre, const std::string& token);
    bool eliminarCola(const std::string& nombre, const std::string& token);
    bool autorizarCola(const std::string& nombre, const std::string& usernameObjetivo, const std::string& token);
    bool enviarMensajeACola(const std::string& nombre, const std::string& contenido, const std::string& token);
    std::optional<Mensaje> consumirMensajeDeCola(const std::string& nombre, const std::string& token);

    // --- Tópicos
    bool crearTopico(const std::string& nombre, const std::string& token);
    bool eliminarTopico(const std::string& nombre, const std::string& token);
    bool suscribirATopico(const std::string& nombre, const std::string& token);
    bool publicarEnTopico(const std::string& nombre, const std::string& contenido, const std::string& token);
    std::vector<Mensaje> consumirDesdeTopico(const std::string& nombre, const std::string& token);

    std::vector<std::string> listarTopicos();
    std::vector<std::string> listarColas();


private:
    Persistencia persistencia;
    std::unordered_map<std::string, std::shared_ptr<Cola>> colas;
    std::unordered_map<std::string, std::shared_ptr<Topico>> topicos;

    bool autenticarYObtenerUsuario(const std::string& token, Usuario& usuario);
};

#endif // BROKER_HPP
