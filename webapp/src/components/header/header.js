import './header.css'

import {Link} from 'react-router-dom';
import React from 'react';

const Header = () => {
    return(
        <div className="bg-purple-900 w-full px-4 py-2 flex justify-center">
            <div className="max-w-3xl w-full ">
                <div className="flex items-center justify-between text-purple-200">
                    <Link to='/'>
                        <div className="text-purple-100  text-2xl font-mono hover:text-purple-300 cursor-pointer select-none">EasySubmit</div>
                    </Link>
                    <div className="flex">
                        <Link to="/signin">
                            <button type="button" className="focus:outline-none active:border-purple-900 hover:text-purple-300 hover:border-purple-300 w-16 h-8 border-2 border-purple-200 rounded-lg  text-xs mr-2 font-thin p-1 ">Sign in</button>
                        </Link>
                        <Link to='/signup'>
                            <button type="button" className="focus:outline-none w-16 h-8 border-2 hover:text-purple-300 hover:border-purple-300 border-purple-200 rounded-lg  text-xs font-thin p-1">Sign up</button>                    
                        </Link>
                    </div>
                </div>
            </div>
        </div>
    )
}

export default Header;