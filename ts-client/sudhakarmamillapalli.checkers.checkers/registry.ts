import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgCreateGame } from "./types/checkers/checkers/tx";
import { MsgPlayMove } from "./types/checkers/checkers/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/sudhakarmamillapalli.checkers.checkers.MsgCreateGame", MsgCreateGame],
    ["/sudhakarmamillapalli.checkers.checkers.MsgPlayMove", MsgPlayMove],
    
];

export { msgTypes }