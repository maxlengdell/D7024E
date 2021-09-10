export localIP=$(hostname -i)

echo Welcome to D7024E cli. Enter your command to interact with the system. 
nc -u $localIP 1010