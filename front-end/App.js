import { NavigationContainer } from '@react-navigation/native'
import Toast from 'react-native-toast-message';

import Main_routs from "./src/stacks/Main_routs"

export default function App() {
  return (
    <NavigationContainer>
        <Main_routs/>
        <Toast/>
    </NavigationContainer>
  );
}