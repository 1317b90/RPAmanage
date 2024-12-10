
export interface taskI {
    ID: number;
    CreatedAt: string;
    UpdatedAt: string;
    DeletedAt: string | null;
    RPAName: string;
    Input: string;
    Output: string;
    State: string;
}

export interface filterI {
    text: string;
    value: string;
}



export interface varI {
    ID: number | null;
    RPAGroup: string;
    RPAName: string;
    VarName: string;
    VarRemark: string;
    AsName: string;
    VarType: string;
    VerifyType: string;
    Default:string
    Required:boolean;
}

export interface RPAGroupI {
    ID: number;
    Name: string;
    Remark: string;
    IP: string;
}

export interface RPAI {
    ID: number;
    Name: string
    Remark: string
    Group: string
    Now:boolean
    Spont:boolean
}