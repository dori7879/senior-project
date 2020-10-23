import { combineReducers, createStore } from "redux";

import { reducer as formReducer } from "redux-form/immutable";

const rootReducer = combineReducers({
  form: formReducer
});

export default createStore(rootReducer);