let Binary = require("./binary");

let schemaJson = `{
  "name": "payload",
  "type": "object",
  "schemas": [
    {"name": "connection", "type": "byte", "le": true},
    {"name": "type", "type": "byte"},
    {"name": "battery", "type": "byte"},
    {
      "name": "cell",
      "type": "[]object",
      "schemas": [
        {"name": "cellId", "type": "[4]byte"},
        {"name": "lac", "type": "int16"},
        {"name": "mcc", "type": "int16"},
        {"name": "mnc", "type": "int16"},
        {"name": "sig", "type": "byte"}
      ]
    },
    {
      "name": "wifi",
      "type": "[]object",
      "schemas": [
        {"name": "mac", "type": "[6]byte"},
        {"name": "sig", "type": "byte"}
      ]
    },
    {"name": "test", "type": "[]int16"}
  ]
}`;

let bytes = new Uint8Array([
  0x01, 0x01, 0x63, 0x02, 0xce, 0x11, 0x1d, 0xc1, 0x1a, 0x00, 0xc1, 0x0c, 0x10, 0x1c, 0x51, 0xce, 0x11, 0x1d, 0xc2,
  0x1a, 0x00, 0xc2, 0x0c, 0x20, 0x2c, 0x52, 0x03, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x0f, 0xa1, 0xb2, 0xc3, 0xd4,
  0xe5, 0xf6, 0x0f, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x0f, 0x03, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06
]);

function test(bytes, schemaJson) {
  console.log(`length = ${bytes.length}`);
  console.time("all");
  let payload = new Binary(schemaJson, bytes);
  let object = payload.toObject();
  console.timeLog("all");

  console.time("log");
  console.log(`----------------------`);
  console.log(JSON.stringify(object, null, 2));
  console.log(`----------------------`);
  console.log(payload.isDone());
  console.timeEnd("log");
  return payload;
}

// test(bytes, schemaJson);

let schemaJson2 = {
  "name": "payload",
  "type": "object",
  "schemas": [
    {"name": "id", "type": "int32"},
    {"name": "ts", "type": "int32"},
    {"name": "status", "type": "int8"},
    {"name": "total_vol", "type": "int32"},
    {"name": "positive_vol", "type": "int32"},
    {"name": "negative_vol", "type": "int32"},
    {"name": "current_flow_rate", "type": "int32"},
    {"name": "current_water_temp", "type": "int32"},
    {"name": "error_code", "type": "int16"},
    {"name": "oldest_log_ts", "type": "int32"},
    {"name": "oldest_log_vol", "type": "int32"},
    {"name": "delta_vol", "type": "[47]int16"},
    {"name": "end_of_data", "type": "int16"},
    {"name": "TX_time", "type": "uint32"},
    {"name": "RX_time", "type": "uint32"},
    {"name": "last_TX_time", "type": "uint32"},
    {"name": "last_RX_time", "type": "uint32"},
    {"name": "TX_power", "type": "int16"},
    {"name": "CSQ", "type": "uint16"},
    {"name": "TX_fail_count", "type": "uint16"},
    {"name": "RSRQ", "type": "int16"},
    {"name": "OPERATOR_MODE", "type": "uint8"},
    {"name": "CURRENT_BAND", "type": "uint8"},
    {"name": "ACTIVE_MODE_LAST_DURATION", "type": "uint32"},
    {"name": "sent_to_server", "type": "uint16"},
    {"name": "server_responded_ok", "type": "uint16"},
    {"name": "PSM_fail_count", "type": "uint16"},
    {"name": "Reg_fail_count", "type": "uint16"},
    {"name": "Con_fail_count", "type": "uint16"},
    {"name": "modem_cycle_duration", "type": "uint16"},
    {"name": "IMEI", "type": "[15]int8"}
  ]
};
let bytes2 = new Uint8Array([
  0x01, 0x00, 0x40, 0xd0, 0x00, 0x43, 0x52, 0x76, 0x05, 0x2e, 0x74, 0x78, 0x62, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x13, 0x6b, 0x6d, 0x62, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xf0, 0xd8, 0xff, 0xff, 0x01, 0x00, 0xa0, 0xd2, 0x75, 0x62, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xaa, 0xaa, 0x17, 0x34, 0x00, 0x00, 0x88, 0xee, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xe6, 0x00, 0x08, 0x00, 0x00, 0x00, 0x6a, 0xff, 0x03, 0x03, 0x99, 0x48, 0x00, 0x00, 0x14, 0x00, 0x00, 0x00, 0x14, 0x00, 0x0a, 0x00, 0x2f, 0x00, 0x4a, 0x00, 0x01, 0x32, 0xfc, 0x8b, 0xa7, 0xe2, 0xff, 0x6e, 0x09, 0x5c, 0x00, 0xab, 0xfb, 0xf1, 0xff, 0x61, 0x7d, 0xe5
]);

let payload = test(bytes2, schemaJson2);
console.log();