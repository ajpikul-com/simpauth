export type signature = {
    Format: string;
    Blob: string;
    Rest: string | null;
};
export type cookie = {
    StateString: string;
    Sig: signature;
};
export declare function getCookie(name: string): string;
export declare function decodeCookie(cookie: string): cookie | null;
//# sourceMappingURL=index.d.ts.map