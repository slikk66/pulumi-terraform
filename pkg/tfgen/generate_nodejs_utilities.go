package tfgen

const tsUtilitiesFile = `
import * as pulumi from "@pulumi/pulumi";

export function getEnv(...vars: string[]): string | undefined {
    for (const v of vars) {
        const value = process.env[v];
        if (value) {
            return value;
        }
    }
    return undefined;
}

export function getEnvBoolean(...vars: string[]): boolean | undefined {
    const s = getEnv(...vars);
    if (s !== undefined) {
        // NOTE: these values are taken from https://golang.org/src/strconv/atob.go?s=351:391#L1, which is what
        // Terraform uses internally when parsing boolean values.
        if (["1", "t", "T", "true", "TRUE", "True"].find(v => v === s) !== undefined) {
            return true;
        }
        if (["0", "f", "F", "false", "FALSE", "False"].find(v => v === s) !== undefined) {
            return false;
        }
    }
    return undefined;
}

export function getEnvNumber(...vars: string[]): number | undefined {
    const s = getEnv(...vars);
    if (s !== undefined) {
        const f = parseFloat(s);
        if (!isNaN(f)) {
            return f;
        }
    }
    return undefined;
}

export function requireWithDefault<T>(req: () => T, def: T | undefined): T {
    try {
        return req();
    } catch (err) {
        if (def === undefined) {
            throw err;
        }
    }
    return def;
}

export function unwrap(val: pulumi.Input<any>): pulumi.Output<any> {
    if (val === null || typeof val !== "object") {
        return pulumi.output(val);
    } else if (val instanceof Promise) {
        return pulumi.output(val).apply(unwrap);
    } else if (pulumi.Output.isInstance(val)) {
        return val.apply(unwrap);
    } else if (val instanceof Array) {
        return pulumi.all(val.map(unwrap));
    } else {
        const unwrappedObject: any = {};
        Object.keys(val).forEach(k => {
            unwrappedObject[k] = unwrap(val[k]);
        });

        return pulumi.all(unwrappedObject);
    }
}
`
