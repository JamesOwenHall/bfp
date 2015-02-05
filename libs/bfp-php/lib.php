<?php

class BFP {
  private $type;
  private $addr;

  public function __construct($type, $addr) {
    $this->type = $type;
    $this->addr = $addr;
  }

  public function getType() {
    return $this->type;
  }

  public function getAddr() {
    return $this->addr;
  }

  public function hit($direction, $value) {
    $file = stream_socket_client($this->type . '://' . $this->addr);
    if ($file === FALSE) {
      return TRUE;
    }

    $request = json_encode(array(
      'Direction' => $direction,
      'Value' => $value,
    ));
    fwrite($file, $request);

    $responseJson = fgets($file);
    $response = json_decode($responseJson, TRUE);

    if ($response['Valid'] === FALSE) {
      return FALSE;
    } else {
      return TRUE;
    }
  }
}
