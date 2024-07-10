# BAIA - Business AI Assistant for WhatsApp

BAIA es un bot de inteligencia artificial diseñado para asistir a negocios en la plataforma de mensajería WhatsApp. Este proyecto utiliza el poder de la AI para interactuar con los clientes, responder preguntas frecuentes, procesar órdenes y mucho más.

## Características

- **Interacción Automática**: Responde automáticamente a los mensajes de los clientes.
- **Procesamiento de Lenguaje Natural (NLP)**: Comprende y procesa preguntas y comandos en lenguaje natural.
- **Gestión de Órdenes**: Ayuda en la recepción y gestión de órdenes de productos o servicios.
- **Personalización**: Personaliza las respuestas basadas en el historial del cliente y preferencias.
- **Integración Fácil**: Se integra fácilmente con la API de WhatsApp Business.

## Requisitos

- Python 3.7 o superior
- Flask
- TensorFlow o PyTorch (para modelos de AI)
- Twilio API para WhatsApp

## Instalación

1. Clona el repositorio:
    ```bash
    git clone https://github.com/galliard-ai/go_baia.git
    cd go_baia
    ```

2. Crea un entorno virtual:
    ```bash
    python3 -m venv venv
    source venv/bin/activate  # En Windows usa `venv\Scripts\activate`
    ```

3. Instala las dependencias:
    ```bash
    pip install -r requirements.txt
    ```

4. Configura las variables de entorno necesarias en un archivo `.env`:
    ```env
    TWILIO_ACCOUNT_SID=your_account_sid
    TWILIO_AUTH_TOKEN=your_auth_token
    WHATSAPP_NUMBER=your_whatsapp_number
    ```

5. Ejecuta la aplicación:
    ```bash
    flask run
    ```

## Uso

Una vez que la aplicación está en funcionamiento, estará lista para recibir y responder mensajes en WhatsApp. Para probar el bot:

1. Envía un mensaje al número de WhatsApp configurado.
2. El bot procesará el mensaje y responderá automáticamente.

## Contribución

Contribuciones son bienvenidas. Para contribuir:

1. Haz un fork del proyecto.
2. Crea una nueva rama (`git checkout -b feature/nueva-funcionalidad`).
3. Realiza tus cambios y haz commit (`git commit -am 'Agrega nueva funcionalidad'`).
4. Haz push a la rama (`git push origin feature/nueva-funcionalidad`).
5. Abre un Pull Request.

## Licencia

Este proyecto está bajo la licencia MIT. Para más detalles, ver el archivo [LICENSE](LICENSE).

## Contacto

Para cualquier consulta o soporte, puedes abrir un issue en el repositorio o contactar al equipo de desarrollo en [correo@dominio.com](mailto:correo@dominio.com).

---

¡Gracias por usar BAIA! Tu asistente de negocios inteligente en WhatsApp.
