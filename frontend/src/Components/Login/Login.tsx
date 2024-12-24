import React from "react";
import { Button } from "@mui/material";
import "./Login.scss";
import { ButtonText } from '../../Constants/text';
import { urls } from "../../Constants/urls";

export const Login = () => {

    const login = async () => {
        try {
            window.location.href = urls.login;
        } catch (error) {
            console.error("Login Failed: ", error);
        }
    }

    return (
        <div className="login">
            <Button 
                className="login-button" 
                variant="contained" 
                color="primary"
                onClick={login}
            >
                {ButtonText.loginYahoo}
            </Button>
        </div>
    )
}

export default Login