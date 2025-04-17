#ifndef USUARIO_HPP
#define USUARIO_HPP

#include <string>
#include <jwt-cpp/jwt.h>

class Usuario {
public:
    Usuario(const std::string& username);

    std::string getUsername() const;
    std::string generarToken();
    static bool verificarToken(const std::string& token, std::string& username_out);
    std::chrono::system_clock::time_point getExpiracion() const;


private:
    std::string username;
    std::chrono::system_clock::time_point expiracion;
    inline static const std::string SECRET = "clave_super_secreta";
};

#endif // USUARIO_HPP
