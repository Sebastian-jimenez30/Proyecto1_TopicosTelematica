#ifndef QUEUE_HPP
#define QUEUE_HPP

#include <string>
#include <queue>
#include "message.hpp"
#include "persistence.hpp"
#include "user.hpp"

class Cola {
public:
    Cola(const std::string& nombre, Usuario& creador, Persistencia& persistencia);

    const std::string& getNombre() const;
    const std::string& getCreadorUsername() const;

    bool encolar(const Mensaje& mensaje);     
    bool hayMensajes() const;
    Mensaje desencolar();

    bool puedeEliminar(const Usuario& usuario) const;

    bool eliminar(const Usuario& usuario);
    bool autorizarUsuario(const std::string& usernameAutorizado);
    bool usuarioPuedeConsumir(const Usuario& usuario) const;


private:
    std::string nombre;
    std::string creadorUsername;
    std::queue<Mensaje> mensajes;
    Persistencia& db;

    void cargarMensajesDesdeBD();
};

#endif
