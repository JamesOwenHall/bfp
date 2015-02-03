<?php

class BFP {
  private $addr;

  public function __construct($addr) {
    $this->addr = $addr;
  }

  public function getAddr() {
    return $this->addr;
  }

  public function hit($direction, $value) {
    $file = stream_socket_client('unix://' . $this->addr);
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
