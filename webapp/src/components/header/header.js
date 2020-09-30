import './header.css'

import React from 'react';

const Header = () => {
    return(
        <div className="bg-purple-600 w-full px-4 py-2 flex justify-center">
            <div className="max-w-3xl w-full ">
                <div className="flex items-center justify-between text-purple-200">
                    <div className="text-purple-900 font-bold text-lg font-mono hover:text-purple-800 cursor-pointer select-none">EasySubmit</div>
                    <div className="flex">
                        <button type="button" className="border-2 border-purple-200 rounded-lg  text-xs mr-2 font-thin p-1 ">Sign in</button>
                        <button type="button" className="border-2 border-purple-200 rounded-lg  text-xs font-thin p-1">Sign up</button>
                    </div>
                </div>
            </div>
        </div>
    )
}

export default Header;