import { NavigationContainer } from '@react-navigation/native'
import Toast from 'react-native-toast-message';

import MainStack from "./src/navigation/MainStack"

export default function App() {
  return (
    <NavigationContainer>
        <MainStack/>
        <Toast/>
    </NavigationContainer>
  );
}