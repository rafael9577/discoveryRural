import { View, Text, TouchableOpacity, TextInput, Image, ActivityIndicator, Modal } from "react-native";
import { SafeAreaView, } from "react-native-safe-area-context";
import { useNavigation } from "@react-navigation/native";
import Toast from "react-native-toast-message";
import { Feather } from "@expo/vector-icons"
import { useState } from 'react';
import styled from "./styled";


export default () => {

    const [emailFild, setEmailFild] = 'insira o email';
    const [passowrdFild, setPassowrdFild] = 'insira a Senha';
    const [passowrdFild1, setPassowrdFild1] = 'confirme a senha';

    <SafeAreaView>
        <View>

            <TextInput
                placeholder="insira seu e-mail"
                onChangeText={t => setEmailFild(t)}
            />
            <TextInput
                placeholder="insira a senha"
                onChangeText={t => setPassowrdFild(t)}
            />
            <TextInput
                placeholder="confirme a senha"
                onChangeText={t => setPassowrdFild1(t)}
            />

        </View>

        <View>

        </View>
    </SafeAreaView>
}
