#!/usr/bin/env python3

import requests
import json
import getpass
from colorama import Fore, Style, init

# Inicializar colorama
init(autoreset=True)

BASE_URL = "http://192.168.1.60:8080"

def menu():
    print(Fore.CYAN + "\n=== Menú Principal ===")
    print("1. Registrar usuario")
    print("2. Iniciar sesión")
    print("3. Crear cola")
    print("4. Eliminar cola")
    print("5. Autorizar usuario en cola")
    print("6. Enviar mensaje a cola")
    print("7. Consumir mensaje de cola")
    print("8. Listar colas")
    print("9. Crear tópico")
    print("10. Eliminar tópico")
    print("11. Suscribirse a tópico")
    print("12. Publicar mensaje en tópico")
    print("13. Consumir mensajes de tópico")
    print("14. Listar tópicos")
    print("15. Salir")
    return input(Fore.GREEN + "Selecciona una opción: ")

def registrar_usuario():
    print(Fore.CYAN + "\n=== Registrar Usuario ===")
    username = input("Nombre de usuario: ")
    password = getpass.getpass("Contraseña: ")
    data = {"username": username, "password": password}
    response = requests.post(f"{BASE_URL}/register", json=data)
    print_response(response)

def iniciar_sesion():
    global token
    print(Fore.CYAN + "\n=== Iniciar Sesión ===")
    username = input("Nombre de usuario: ")
    password = getpass.getpass("Contraseña: ")
    data = {"username": username, "password": password}
    try:
        response = requests.post(f"{BASE_URL}/login", json=data)
        if response.status_code == 200:
            token = response.json().get("token")
            print(Fore.GREEN + "Inicio de sesión exitoso. Token guardado.")
        else:
            print_response(response)
    except requests.exceptions.ConnectionError:
        print(Fore.RED + "Error: No se pudo conectar al servidor. Verifica que esté en ejecución.")

def crear_cola():
    print(Fore.CYAN + "\n=== Crear Cola ===")
    nombre = input("Nombre de la cola: ")
    headers = {"Authorization": f"Bearer {token}"}
    data = {"nombre": nombre}
    response = requests.post(f"{BASE_URL}/colas", json=data, headers=headers)
    print_response(response)

def eliminar_cola():
    print(Fore.CYAN + "\n=== Eliminar Cola ===")
    nombre = input("Nombre de la cola a eliminar: ")
    headers = {"Authorization": f"Bearer {token}"}
    response = requests.delete(f"{BASE_URL}/colas/{nombre}", headers=headers)
    print_response(response)

def autorizar_usuario():
    print(Fore.CYAN + "\n=== Autorizar Usuario en Cola ===")
    nombre = input("Nombre de la cola: ")
    usuario = input("Usuario a autorizar: ")
    headers = {"Authorization": f"Bearer {token}"}
    data = {"usuario": usuario}
    response = requests.post(f"{BASE_URL}/colas/{nombre}/autorizar", json=data, headers=headers)
    print_response(response)

def enviar_mensaje():
    print(Fore.CYAN + "\n=== Enviar Mensaje a Cola ===")
    nombre = input("Nombre de la cola: ")
    contenido = input("Contenido del mensaje: ")
    headers = {"Authorization": f"Bearer {token}"}
    data = {"contenido": contenido}
    response = requests.post(f"{BASE_URL}/colas/{nombre}/enviar", json=data, headers=headers)
    print_response(response)

def consumir_mensaje():
    print(Fore.CYAN + "\n=== Consumir Mensaje de Cola ===")
    nombre = input("Nombre de la cola: ")
    headers = {"Authorization": f"Bearer {token}"}
    response = requests.get(f"{BASE_URL}/colas/{nombre}/consumir", headers=headers)
    print_response(response)

def listar_colas():
    print(Fore.CYAN + "\n=== Listar Colas ===")
    headers = {"Authorization": f"Bearer {token}"}
    response = requests.get(f"{BASE_URL}/colas", headers=headers)
    print_response(response)

def crear_topico():
    print(Fore.CYAN + "\n=== Crear Tópico ===")
    nombre = input("Nombre del tópico: ")
    headers = {"Authorization": f"Bearer {token}"}
    data = {"nombre": nombre}
    response = requests.post(f"{BASE_URL}/topicos", json=data, headers=headers)
    print_response(response)

def eliminar_topico():
    print(Fore.CYAN + "\n=== Eliminar Tópico ===")
    nombre = input("Nombre del tópico a eliminar: ")
    headers = {"Authorization": f"Bearer {token}"}
    response = requests.delete(f"{BASE_URL}/topicos/{nombre}", headers=headers)
    print_response(response)

def suscribirse_topico():
    print(Fore.CYAN + "\n=== Suscribirse a Tópico ===")
    nombre = input("Nombre del tópico: ")
    headers = {"Authorization": f"Bearer {token}"}
    data = {"nombre": nombre}
    response = requests.post(f"{BASE_URL}/topicos/{nombre}/suscribir", json=data, headers=headers)
    print_response(response)

def publicar_topico():
    print(Fore.CYAN + "\n=== Publicar Mensaje en Tópico ===")
    nombre = input("Nombre del tópico: ")
    contenido = input("Contenido del mensaje: ")
    headers = {"Authorization": f"Bearer {token}"}
    data = {"contenido": contenido}
    response = requests.post(f"{BASE_URL}/topicos/{nombre}/publicar", json=data, headers=headers)
    print_response(response)

def consumir_topico():
    print(Fore.CYAN + "\n=== Consumir Mensajes de Tópico ===")
    nombre = input("Nombre del tópico: ")
    headers = {"Authorization": f"Bearer {token}"}
    response = requests.get(f"{BASE_URL}/topicos/{nombre}/consumir", headers=headers)
    print_response(response)

def listar_topicos():
    print(Fore.CYAN + "\n=== Listar Tópicos ===")
    headers = {"Authorization": f"Bearer {token}"}
    response = requests.get(f"{BASE_URL}/topicos", headers=headers)
    print_response(response)

def print_response(response):
    try:
        print(Fore.GREEN + "\nRespuesta del servidor:")
        print(json.dumps(response.json(), indent=4))
    except json.JSONDecodeError:
        print(Fore.RED + response.text)

if __name__ == "__main__":
    token = None
    while True:
        opcion = menu()
        if opcion == "1":
            registrar_usuario()
        elif opcion == "2":
            iniciar_sesion()
        elif opcion == "3":
            if token:
                crear_cola()
            else:
                print(Fore.RED + "Por favor, inicia sesión primero.")
        elif opcion == "4":
            if token:
                eliminar_cola()
            else:
                print(Fore.RED + "Por favor, inicia sesión primero.")
        elif opcion == "5":
            if token:
                autorizar_usuario()
            else:
                print(Fore.RED + "Por favor, inicia sesión primero.")
        elif opcion == "6":
            if token:
                enviar_mensaje()
            else:
                print(Fore.RED + "Por favor, inicia sesión primero.")
        elif opcion == "7":
            if token:
                consumir_mensaje()
            else:
                print(Fore.RED + "Por favor, inicia sesión primero.")
        elif opcion == "8":
            if token:
                listar_colas()
            else:
                print(Fore.RED + "Por favor, inicia sesión primero.")
        elif opcion == "9":
            if token:
                crear_topico()
            else:
                print(Fore.RED + "Por favor, inicia sesión primero.")
        elif opcion == "10":
            if token:
                eliminar_topico()
            else:
                print(Fore.RED + "Por favor, inicia sesión primero.")
        elif opcion == "11":
            if token:
                suscribirse_topico()
            else:
                print(Fore.RED + "Por favor, inicia sesión primero.")
        elif opcion == "12":
            if token:
                publicar_topico()
            else:
                print(Fore.RED + "Por favor, inicia sesión primero.")
        elif opcion == "13":
            if token:
                consumir_topico()
            else:
                print(Fore.RED + "Por favor, inicia sesión primero.")
        elif opcion == "14":
            if token:
                listar_topicos()
            else:
                print(Fore.RED + "Por favor, inicia sesión primero.")
        elif opcion == "15":
            print(Fore.YELLOW + "Saliendo del cliente...")
            break
        else:
            print(Fore.RED + "Opción no válida. Intenta de nuevo.")