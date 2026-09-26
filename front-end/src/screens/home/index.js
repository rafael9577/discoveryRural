import { SafeAreaView, } from "react-native-safe-area-context";
import { Text, View, TouchableOpacity } from "react-native";
import { useNavigation } from "@react-navigation/native";
import Styled from "./style";

export default () => {

    const Navigator = useNavigation();

    function sair() {
        Navigator.reset({
            routes: [{ name: 'SingIn' }]
        })
    }

    return (
        <SafeAreaView style={Styled.safeArea}>

           
        </SafeAreaView>
    );
};