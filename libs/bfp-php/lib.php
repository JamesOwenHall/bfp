<?php

class BFP {
  private $type;
  private $addr;
  private $port;

  public function __construct($type, $addr, $port = 9999) {
    $this->type = $type;
    $this->addr = $addr;
    $this->port = $port;
  }

  public function getType() {
    return $this->type;
  }

  public function getAddr() {
    return $this->addr;
  }

  public function hit($direction, $value) {
    $file = NULL;
    if ($this->type === "unix") {
      $file = pfsockopen("unix://" . $this->addr);
    } else {
      $file = pfsockopen("tcp://" . $this->addr, $this->port);
    }

    if ($file === FALSE) {
      return FALSE;
    }

    $request = json_encode(array(
      'Direction' => $direction,
      'Value' => $value,
    ));
    fwrite($file, $request);

    $valid = $this->readResponse($file);
    return $valid;
  }

  private function readResponse($file) {
    $data = fread($file, 1);
    if ($data === "t") {
      return TRUE;
    } else {
      return FALSE;
    }
  }
}
