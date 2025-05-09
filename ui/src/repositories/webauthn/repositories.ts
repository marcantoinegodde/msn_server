import axios from "axios";

export const postWebauthnRegisterBegin = async () => {
  const response = await axios.post<PublicKeyCredentialCreationOptionsJSON>(
    `${import.meta.env.VITE_API_URL}/webauthn/register/begin`,
    {},
    { withCredentials: true }
  );
  return response.data;
};

export const postWebauthnRegisterFinish = async (body: PublicKeyCredential) => {
  const response = await axios.post<PublicKeyCredential>(
    `${import.meta.env.VITE_API_URL}/webauthn/register/finish`,
    body,
    { withCredentials: true }
  );
  return response.data;
};

export const postWebauthnLoginBegin = async () => {
  const response = await axios.post<PublicKeyCredentialCreationOptionsJSON>(
    `${import.meta.env.VITE_API_URL}/webauthn/login/begin`,
    {},
    { withCredentials: true }
  );
  return response.data;
};

export const postWebauthnLoginFinish = async (body: PublicKeyCredential) => {
  const response = await axios.post<PublicKeyCredential>(
    `${import.meta.env.VITE_API_URL}/webauthn/login/finish`,
    body,
    { withCredentials: true }
  );
  return response.data;
};
