DELIMITER $$

CREATE FUNCTION obtener_documentos_usuario(uid INT)
RETURNS JSON
DETERMINISTIC
BEGIN
  DECLARE cartas JSON;
  DECLARE escrituras JSON;

  SET cartas = (
    SELECT JSON_ARRAYAGG(
      JSON_OBJECT(
        'id_carta_poder', ID_CARTA_PODER,
        'poderante', NOMBRE_PODERANTE,
        'apoderado', NOMBRE_APODERADO,
        'domicilio', DOMICILIO,
        'telefono', TELEFONO,
        'rfc', RFC,
        'id_pago', ID_PAGO
      )
    )
    FROM CARTA_PODER
    WHERE ID_USUARIO = uid
  );

  SET escrituras = (
    SELECT JSON_ARRAYAGG(
      JSON_OBJECT(
        'id_escritura_publica', E.ID_ESCRITURA_PUBLICA,
        'nombre', E.NOMBRE,
        'nacionalidad', E.NACIONALIDAD,
        'estado_civil', C.NOMBRE,
        'ocupacion', E.OCUPACION,
        'curp', E.CURP,
        'id_pago', E.ID_PAGO
      )
    )
    FROM ESCRITURA_PUBLICA E
    JOIN ESTADOS_CIVILES C ON E.ID_ESTADO_CIVIL = C.ID_ESTADO_CIVIL
    WHERE E.ID_USUARIO = uid
  );

  RETURN JSON_OBJECT(
    'cartas_poder', IFNULL(cartas, JSON_ARRAY()),
    'escrituras_publicas', IFNULL(escrituras, JSON_ARRAY())
  );
END $$

DELIMITER ;
