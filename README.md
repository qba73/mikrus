# mikrus

Go client library for [MIKRUS VPS](https://mikr.us) Provider

## Endpoints

```shell
/info       - informacje o Twoim serwerze (cache=60s)
/serwery    - listuje wszystkie Twoje serwery (cache=60s)
/restart    - restartuje Twój serwer
/logs       - podgląd ostatnich logów [10 sztuk]
/logs/ID    - podgląd konkretnego wpisu w logach (po ID)
/amfetamina - uruchamia amfetaminę na serwerze (zwiększenie parametrów)
/db         - zwraca dane dostępowe do baz danych (cache=60s);
/exec       - wywołuje polecenie/polecenia wysłane w zmiennej 'cmd' (POST)
/stats      - statystyki użycia dysku, pamięci, uptime itp. (cache=60s)
/porty      - zwraca przypisane do Twojego serwera porty TCP/UDP (cache=60s)
/cloud      - zwraca listę usług cloud przypisanych do Twojego konta wraz ze statystykami
/domain     - podaj port w zmiennej port i domenę w zmiennej domena lub `-` by system sam nadał ci domenę
```
