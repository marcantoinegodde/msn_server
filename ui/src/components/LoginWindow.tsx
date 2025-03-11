import keys from "../icons/keys.png";

function LoginWindows() {
  return (
    <div className="window w-[400px]">
      <div className="title-bar">
        <div className="title-bar-text">Log On - MSN Messenger Service</div>
      </div>
      <div className="window-body">
        <div className="flex space-x-2.5">
          <div>
            <img src={keys} />
          </div>
          <div className="w-full">
            <p>Enter your e-mail address and password:</p>
            <fieldset>
              <div className="field-row">
                <label className="text-nowrap">Logon Name:</label>
                <input type="email" className="w-full" />
              </div>
              <div className="field-row">
                <label>Password:</label>
                <input type="password" className="w-full" />
              </div>
            </fieldset>
          </div>
        </div>
        <div className="flex justify-end mt-2.5">
          <input type="submit" />
        </div>
      </div>
    </div>
  );
}

export default LoginWindows;
