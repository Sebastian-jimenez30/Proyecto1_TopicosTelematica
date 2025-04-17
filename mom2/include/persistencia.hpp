#ifndef PERSISTENCIA_HPP
#define PERSISTENCIA_HPP

#include <sqlite3.h>
#include <string>
#include <vector>
#include "mensaje.hpp"

class Persistencia {
public:
    Persistencia(const std::string& db_path);
    ~Persistencia();

    bool inicializarBaseDeDatos();

    // Usuarios
    bool crearUsuario(const std::string& username, const std::string& passwordHash);
    bool verificarCredenciales(const std::string& username, const std::string& passwordHash);
    bool guardarToken(const std::string& username, const std::string& token, const std::string& expiracion);
    bool validarTokenActivo(const std::string& token);

    // Colas
    bool crearCola(const std::string& nombreCola, const std::string& creadorUsername);
    bool eliminarCola(const std::string& nombreCola, const std::string& username);
    bool autorizarUsuarioParaConsumir(const std::string& username, const std::string& colaNombre);
    bool verificarAutorizacion(const std::string& username, const std::string& colaNombre);

    // Mensajes
    bool guardarMensaje(const Mensaje& mensaje);
    std::vector<Mensaje> cargarMensajesPorCanal(const std::string& canal, const std::string& tipo);
    std::vector<Mensaje> cargarMensajesDesdeOffset(const std::string& topico, int offset, int limite = 10);

    bool crearTopico(const std::string& nombreTopico, const std::string& creadorUsername);
    bool suscribirUsuarioATopico(const std::string& username, const std::string& topicoNombre);
    int obtenerOffset(const std::string& username, const std::string& topicoNombre);
    bool actualizarOffset(const std::string& username, const std::string& topicoNombre, int nuevoOffset);
    bool eliminarTopico(const std::string& nombreTopico, const std::string& username);
    
    bool registrarLog(const std::string& username, const std::string& actividad);

    std::vector<std::string> obtenerNombresTopicos();
    std::vector<std::string> obtenerNombresColas();



private:
    sqlite3* db;
    bool ejecutar(const std::string& sql);
};

#endif // PERSISTENCIA_HPP
