#ifndef MENSAJE_HPP
#define MENSAJE_HPP

#include <string>
#include <chrono>

class Mensaje {
public:
    // Constructor para crear un mensaje nuevo (antes de guardar)
    Mensaje(const std::string& remitente,
            const std::string& contenido,
            const std::string& canal,
            const std::string& tipo,
            std::chrono::system_clock::time_point timestamp);

    // Constructor para cargar desde base de datos (incluye ID)
    Mensaje(int id,
            const std::string& remitente,
            const std::string& contenido,
            const std::string& canal,
            const std::string& tipo,
            std::chrono::system_clock::time_point timestamp);

    // Getters
    int getId() const;
    std::string getRemitente() const;
    std::string getContenido() const;
    std::string getCanal() const;
    std::string getTipo() const;
    std::chrono::system_clock::time_point getTimestamp() const;

    std::string toString() const;

private:
    int id;  // -1 si aún no se ha guardado
    std::string remitente;
    std::string contenido;
    std::string canal;
    std::string tipo;
    std::chrono::system_clock::time_point timestamp;
};

#endif // MENSAJE_HPP
