#ifndef COLA_HPP
#define COLA_HPP

#include <string>
#include <queue>
#include "mensaje.hpp"
#include "persistencia.hpp"
#include "usuario.hpp"

class Cola {
public:
    Cola(const std::string& nombre, Usuario& creador, Persistencia& persistencia);

    const std::string& getNombre() const;
    const std::string& getCreadorUsername() const;

    bool encolar(const Mensaje& mensaje);         // Enviar mensaje
    bool hayMensajes() const;
    Mensaje desencolar();                         // Consumir mensaje

    bool puedeEliminar(const Usuario& usuario) const;

    bool eliminar(const Usuario& usuario);
    bool autorizarUsuario(const std::string& usernameAutorizado);
    bool usuarioPuedeConsumir(const Usuario& usuario) const;


private:
    std::string nombre;
    std::string creadorUsername;
    std::queue<Mensaje> mensajes;
    Persistencia& db;

    void cargarMensajesDesdeBD();                 // Inicializa mensajes desde DB
};

#endif // COLA_HPP
