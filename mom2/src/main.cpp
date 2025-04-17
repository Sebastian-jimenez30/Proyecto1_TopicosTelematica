#include "persistencia.hpp"
#include "usuario.hpp"
#include "topico.hpp"
#include "mensaje.hpp"
#include <iostream>
#include <chrono>

int main() {
    Persistencia db("../data/mom.db");

    if (!db.inicializarBaseDeDatos()) {
        std::cerr << "❌ No se pudo inicializar la base de datos.\n";
        return 1;
    }

    // ----------- Registro de usuarios ------------
    std::string u1 = "sebastian";
    std::string u2 = "andres";

    db.crearUsuario(u1, "1234");
    db.crearUsuario(u2, "5678");

    Usuario usuario1(u1);
    Usuario usuario2(u2);

    // ----------- Crear tópico ------------
    std::string topicoNombre = "eventos";
    Topico topico(topicoNombre, usuario1, db);
    std::cout << "✅ Tópico '" << topicoNombre << "' creado por " << u1 << "\n";

    // ----------- Suscripción ------------
    if (topico.suscribir(usuario2)) {
        std::cout << "📌 Usuario '" << u2 << "' suscrito a '" << topicoNombre << "'.\n";
    }

    // ----------- Publicación de mensajes ------------
    Mensaje m1(u1, "¡Bienvenidos a los eventos!", topicoNombre, "topico", std::chrono::system_clock::now());
    Mensaje m2(u1, "Recuerden registrarse antes de entrar", topicoNombre, "topico", std::chrono::system_clock::now());

    topico.publicarMensaje(m1);
    topico.publicarMensaje(m2);

    std::cout << "📢 Mensajes publicados por " << u1 << " en el tópico '" << topicoNombre << "'\n";

    // ----------- Consumo por parte de andres ------------
    std::vector<Mensaje> recibidos = topico.obtenerMensajesPara(usuario2);

    std::cout << "\n📥 Mensajes recibidos por " << u2 << ":\n";
    for (const auto& msg : recibidos) {
        std::cout << msg.toString() << "\n";
    }

    return 0;
}
