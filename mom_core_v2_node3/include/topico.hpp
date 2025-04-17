#ifndef TOPICO_HPP
#define TOPICO_HPP

#include <string>
#include <vector>
#include "usuario.hpp"
#include "mensaje.hpp"
#include "persistencia.hpp"

class Topico {
public:
    Topico(const std::string& nombre, Usuario& creador, Persistencia& persistencia);

    const std::string& getNombre() const;
    const std::string& getCreadorUsername() const;

    bool publicarMensaje(const Mensaje& mensaje);
    std::vector<Mensaje> obtenerMensajesPara(Usuario& usuario, int maxMensajes = 10);

    bool suscribir(Usuario& usuario);
    bool puedeEliminar(const Usuario& usuario) const;
    bool eliminar(const Usuario& usuario);

private:
    std::string nombre;
    std::string creadorUsername;
    Persistencia& db;
};

#endif // TOPICO_HPP
