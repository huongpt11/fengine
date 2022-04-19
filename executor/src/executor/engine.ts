import {Function, Result, Script, Variable} from "../pb/fengine_pb";
import * as library from "../sdk/db";
import {VM} from "vm2";
import _ from "lodash";

type MsgType = void | number | string | boolean | Uint8Array;
type Func = (input: any) => MsgType;
type ThingReference = {
  [key: string]: MsgType | Func
}

class E {
  private static buildSandbox(script: Script) {
    let code = "";
    const me: ThingReference = {};
    const attributes: ThingReference = {};

    script.getAttributesList().forEach(attr => {
      let name = attr.getName();
      let value: MsgType = E.readVarValue(attr);
      me[name] = value;
      attributes[name] = value;
    });

    script.getRefereeMap().forEach((fn: Function, name: string) => {
      let _code = fn.getCode();
      if (_code) {
        const {params} = E.parseInput(fn.getInputList(), false);
        code += `me['${name}'] = (${params}) => {${_code}};\n`;
      } else {
        me[name] = (input: any) => {
          console.log(input);
        };
      }
    });

    return {sandbox: {me, ...library}, code, attributes};
  }

  private static parseInput(input: Variable[], hasArgs: boolean = true) {
    const args: MsgType[] = [];
    const params: string[] = [];
    input.forEach(inp => {
      params.push(inp.getName());
      if (hasArgs) {
        args.push(E.readVarValue(inp, true));
      } else {
        console.log(`pi: ${inp.getName()}`);
      }
    });
    if (!hasArgs) console.log(`> ${JSON.stringify(args)}`);
    return {args, params};
  }

  private static wrap(output: any, outputType: Variable) {
    const variable = new Variable();
    switch (typeof output) {
      case "object":
        if (outputType.hasJson()) {
          variable.setJson(JSON.stringify(output instanceof Error ? {error: output.message} : output));
        }
        break;
      case "boolean":
        variable.setBool(output);
        break;
      case "number":
        if (Number.isInteger(output)) {
          if (output < 4294967296) {
            variable.setI32(output);
          } else {
            variable.setI64(output);
          }
          break;
        }
        variable.setF64(output);
        break;
      case "string":
        variable.setString(output);
        break;
    }

    return variable;
  }

  private static readVarValue(input: Variable, isParam: boolean = false): MsgType {
    switch (true) {
      case input.hasI32():
        return input.getI32();
      case input.hasI64():
        return input.getI64();
      case input.hasF32():
        return input.getF32();
      case input.hasF64():
        return input.getF64();
      case input.hasBool():
        return input.getBool();
      case input.hasJson():
        return input.getJson();
      case input.hasString():
        return isParam ? `'${input.getString()}'` : input.getString();
      case input.hasBinary():
        return input.getBinary_asU8();
    }
  }

  private static compareAttributes(me: ThingReference, attributes: ThingReference) {
    for (let i in attributes) {
      if (!_.isEqual(attributes[i], me[i])) {
        // if (attributes[i] !== me[i]) {
        console.log(`>>> ${i}: ${attributes[i]} -> ${me[i]}`);
      }
    }
  }

  exec(script: Script): Result {
    try {
      const fn = script.getFunction();
      if (!fn) {
        let json = new Variable().setJson(JSON.stringify({error: "Function is not defined"}));
        return new Result().setOutput(json);
      }

      const {sandbox, code: sandboxCode, attributes} = E.buildSandbox(script);
      const {args, params} = E.parseInput(fn.getInputList());
      const code = `((${params})=>{try{${fn.getCode()}}catch(_e_){return _e_}})(${args.join()})`;
      console.debug(`${JSON.stringify(sandbox)}>---\n${sandboxCode}\n${code}\n---<`);

      const vm = new VM({sandbox});
      const label = new Date().getTime();
      console.time(`${label}`);
      let output = E.wrap(vm.run(sandboxCode + code), fn.getOutput()!);
      console.timeEnd(`${label}`);
      E.compareAttributes(sandbox.me, attributes);

      return new Result().setOutput(output);
    } catch (e: any) {
      return new Result().setOutput(new Variable().setString(e.message));
    }
  }
}

export {E as Executor};
