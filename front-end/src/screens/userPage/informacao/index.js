import { View, Text, TouchableOpacity, TextInput, Image, ActivityIndicator, Modal } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context"

import style from "./style"

export default () => {


    return(
        <SafeAreaView style={style.SafeAreaView}>
            <View style={style.busca}>
                <TextInput style={style.buscaInput} placeholder="aqui"/>
            </View>
        </SafeAreaView>
    )
}