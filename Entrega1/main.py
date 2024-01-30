# RECUPERACIÓ TRAMESA 1
# Autor: Joan Martorell Coll

import threading
import time
import random

SOSPITOSOS = ["Deadshot", "Harley Quinn", "Penguin", "Riddler", "Bane", 
            "Talia al Ghul", "Ra's al Ghul", "Hugo Strange", "Killer Croc", "Catwoman", 
            "Poison Ivy", "Mr. Freeze", "Jason Todd", "Hush", "Joker", 
            "Clayface", "Deathstroke", "Mad Hatter", "Two-Face", "Scarecrow"]
NUM_SOSPITOSOS = 20

# Declaram sospitosos
SOSPITOSOS_DINS_SALA = 1
SOSPITOSOS_FITXATS = 1
SOSPITOSOS_DECLARATS = 1
SOSPITOSOS_ESPERANT = 1

# Declaram variables globals necessàries
FICHAT = True
ACABAT = False

# Declaram semàfors
jutgeSemafor = threading.Semaphore(0)
salaSemafor = threading.Semaphore(1)
fitxarSemafor = threading.Semaphore(1)
esperaJutgeSemafor = threading.Semaphore(0)
esperaSospitosSemafor = threading.Semaphore(0)
declaraSemafor = threading.Semaphore(0)
veredicteSemafor = threading.Semaphore(0)
llibertatSemafor = threading.Semaphore(0)
esperaDeclararSemafor = threading.Semaphore(1)

# Declaram temps de simulació
ESPERA_MIN = 50
ESPERA_MAX = 100
JUTGE_ESPERA_MIN = 100
JUTGE_ESPERA_MAX = 1000

rnd = random.Random()

# Funcions d'espera
def esperaJutge():
    t = rnd.randint(JUTGE_ESPERA_MIN, JUTGE_ESPERA_MAX) / 1000
    time.sleep(t)

def esperaSospitos():
    t = rnd.randint(ESPERA_MIN, ESPERA_MAX) / 1000
    time.sleep(t)

# Funció de Jutge
def jutge():

    global ACABAT, FICHAT, SOSPITOSOS_DINS_SALA
    
    try:
        esperaJutge()  # Espera
        print("----> Jutge Dredd: Jo soc la llei!")
        esperaJutge()
        print("----> Jutge Dredd: Sóc a la sala, tanqueu la porta!")

        jutgeSemafor.release()

        # Si no hi ha ningú s'en va
        if SOSPITOSOS_DINS_SALA == 1:
            print("----> Jutge Dredd: Si no hi ha ningú, me'n vaig!")
            print("----> Jutge Dredd: La justícia descansa, prendré declaració als sospitosos que queden")
            esperaJutgeSemafor.acquire()   
            ACABAT = True
            esperaSospitosSemafor.release(NUM_SOSPITOSOS)
            return

        esperaSospitosSemafor.release(SOSPITOSOS_DINS_SALA - 1) 
        print("----> Jutge Dredd: Fitxeu als sospitosos presents")

        while FICHAT:
            pass

        print("----> Jutge Dredd: Preniu declaració als presents")
        FICHAT = True

        declaraSemafor.release(SOSPITOSOS_DINS_SALA - 1)

        while FICHAT:
            pass

        print("----> Jutge Dredd: Podeu abandonar la sala, tots a l'asil!")
        print("----> Jutge Dredd: La justícia descansa, demà prendré declaració als sospitosos que queden")

        llibertatSemafor.release(SOSPITOSOS_DINS_SALA - 1)
        esperaSospitosSemafor.release(NUM_SOSPITOSOS)

    except Exception as ex:
        print("ERROR: " + str(ex))
    


def sospitos(sospitos):
    
    try:
        global ACABAT, FICHAT, SOSPITOSOS_ESPERANT, SOSPITOSOS_DINS_SALA, SOSPITOSOS_FITXATS, SOSPITOSOS_DECLARATS
        salaSemafor.acquire()

        if jutgeSemafor._value == 0:
            print("      " + sospitos + ": Sóc innocent!")

        esperaSospitos()

        SOSPITOSOS_ESPERANT += 1
        salaSemafor.release()

        if jutgeSemafor._value == 0:
            print("      " + sospitos + " entra al jutjat. Sospitosos: " + str(SOSPITOSOS_DINS_SALA))
            SOSPITOSOS_DINS_SALA += 1

        if SOSPITOSOS_ESPERANT == NUM_SOSPITOSOS:
            esperaJutgeSemafor.release()

        # Esperam fins que el jutge no digui per fitxar
        esperaSospitosSemafor.acquire() 

        if ACABAT:
            print("      " + sospitos + ": No és just, vull declarar! Sóc innocent!")
            return

        # Esperam per fitxar
        fitxarSemafor.acquire()
        print("      " + sospitos + " fitxa. Fitxats: " + str(SOSPITOSOS_FITXATS))
        esperaSospitos()
        SOSPITOSOS_FITXATS += 1
        fitxarSemafor.release()

        if SOSPITOSOS_DINS_SALA == SOSPITOSOS_FITXATS:
            FICHAT = False

        # Esperam per declarar
        declaraSemafor.acquire()
        esperaDeclararSemafor.acquire()
        print("      " + sospitos + " declara. Declaracions: " + str(SOSPITOSOS_DECLARATS))
        SOSPITOSOS_DECLARATS += 1
        esperaDeclararSemafor.release()

        if SOSPITOSOS_DINS_SALA == SOSPITOSOS_DECLARATS:
            FICHAT = False

        llibertatSemafor.acquire()
        ACABAT = True

        print("      " + sospitos + " entra a l'Asil d'Arkham")

    except Exception as ex:
        print("Error: " + str(ex))



if __name__ == "__main__":
    jutge_thread = threading.Thread(target=jutge)
    jutge_thread.start()

    sospitos_threads = []
    for i in range(NUM_SOSPITOSOS):
        sospitos_thread = threading.Thread(target=sospitos, args=(SOSPITOSOS[i],))
        sospitos_threads.append(sospitos_thread)
        sospitos_thread.start()

    jutge_thread.join()
    for sospitos_thread in sospitos_threads:
        sospitos_thread.join()
